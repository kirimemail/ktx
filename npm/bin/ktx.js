#!/usr/bin/env node

// @kirimemail/ktx - self-bootstrapping CLI shim.
// Downloads the correct native binary from GitHub Releases on first run
// and caches it at ~/.local/share/@kirimemail/ktx/<version>/.

const { spawn } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");
const os = require("os");

// ---------------------------------------------------------------------------
// Platform detection
// ---------------------------------------------------------------------------

const platformMap = {
	darwin: "darwin",
	linux: "linux",
	win32: "windows",
};

const archMap = {
	x64: "amd64",
	arm64: "arm64",
};

const osName = platformMap[process.platform];
const archName = archMap[process.arch];
const isWin = process.platform === "win32";

if (!osName || !archName) {
	console.error(
		`ktx is not available for platform ${process.platform}-${process.arch}\n` +
			`Supported: linux (x64, arm64), darwin (x64, arm64), win32 (x64, arm64)`,
	);
	process.exit(1);
}

// ---------------------------------------------------------------------------
// Version & paths
// ---------------------------------------------------------------------------

// package.json sits one level up from bin/
const pkg = require(path.join(__dirname, "..", "package.json"));
const VERSION = pkg.version;

const binaryName = isWin ? "ktx.exe" : "ktx";
const assetName = `ktx-${osName}-${archName}${isWin ? ".exe" : ""}`;

// Binary cache dir — versioned so upgrades re-download
const dataDir =
	process.env.XDG_DATA_HOME || path.join(os.homedir(), ".local", "share");
const cacheDir = path.join(dataDir, "@kirimemail", "ktx", VERSION);
const binaryPath = path.join(cacheDir, binaryName);
const stampPath = path.join(cacheDir, ".download-stamp");

// ---------------------------------------------------------------------------
// Download helpers
// ---------------------------------------------------------------------------

/** Run a single async function with an optional timeout. */
function withTimeout(promise, ms = 30_000) {
	return Promise.race([
		promise,
		new Promise((_, reject) =>
			setTimeout(() => reject(new Error("Download timed out")), ms),
		),
	]);
}

/** Ensure the cache directory exists. */
function ensureDir(dir) {
	fs.mkdirSync(dir, { recursive: true });
}

/** Stream a file from a URL to disk. */
function download(url, dest) {
	return new Promise((resolve, reject) => {
		const file = fs.createWriteStream(dest, { mode: 0o755 });
		const onError = (err) => {
			file.destroy();
			try {
				fs.unlinkSync(dest);
			} catch {
				/* ignore */
			}
			reject(err);
		};

		https
			.get(url, (response) => {
				// Follow redirects (one hop)
				if (
					response.statusCode >= 300 &&
					response.statusCode < 400 &&
					response.headers.location
				) {
					file.close();
					return download(response.headers.location, dest).then(
						resolve,
						reject,
					);
				}

				if (response.statusCode !== 200) {
					return onError(
						new Error(`HTTP ${response.statusCode} downloading ${url}`),
					);
				}

				response.pipe(file);
				file.on("finish", () => file.close(resolve));
			})
			.on("error", onError);
	});
}

/**
 * Download the platform binary from GitHub Releases.
 * Expected URL format:
 *   https://github.com/kirimemail/ktx/releases/download/v{version}/ktx-{os}-{arch}
 */
async function downloadBinary() {
	const url = [
		"https://github.com/kirimemail/ktx/releases/download",
		`v${VERSION}`,
		assetName,
	].join("/");

	console.error(`ktx: downloading binary for ${osName}-${archName}...`);
	console.error(`     ${url}`);

	ensureDir(cacheDir);
	await withTimeout(download(url, binaryPath));

	// Write stamp so we don't re-download
	fs.writeFileSync(stampPath, VERSION, "utf-8");
	console.error("ktx: download complete");
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

(async () => {
	// Check if the binary is already cached
	const binaryMissing = !fs.existsSync(binaryPath);
	const stampMissing = !fs.existsSync(stampPath);
	const versionMismatch =
		!stampMissing && fs.readFileSync(stampPath, "utf-8").trim() !== VERSION;

	if (binaryMissing || stampMissing || versionMismatch) {
		try {
			await downloadBinary();
		} catch (err) {
			console.error(`ktx: failed to download binary: ${err.message}`);
			console.error(
				`ktx: you can manually download ${assetName} from`,
				`https://github.com/kirimemail/ktx/releases/tag/v${VERSION}`,
			);
			process.exit(1);
		}
	}

	// Spawn the binary with the original args
	const child = spawn(binaryPath, process.argv.slice(2), {
		stdio: "inherit",
		windowsHide: true,
	});

	child.on("exit", (code, signal) => {
		if (code !== null) process.exit(code);
		// On signal, forward it
		process.kill(process.pid, signal);
	});

	child.on("error", (err) => {
		console.error(`Failed to spawn ktx: ${err.message}`);
		console.error(`Expected at: ${binaryPath}`);
		process.exit(1);
	});
})();
