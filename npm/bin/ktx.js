#!/usr/bin/env node

// @kirimemail/ktx - JS bin shim that spawns the correct native binary
// for the current platform.

const { spawn } = require("child_process");
const path = require("path");

const knownPlatforms = {
	"linux-x64": "@kirimemail/ktx-linux-x64",
	"linux-arm64": "@kirimemail/ktx-linux-arm64",
	"darwin-x64": "@kirimemail/ktx-darwin-x64",
	"darwin-arm64": "@kirimemail/ktx-darwin-arm64",
	"win32-x64": "@kirimemail/ktx-win32-x64",
	"win32-arm64": "@kirimemail/ktx-win32-arm64",
};

const key = `${process.platform}-${process.arch}`;
const pkg = knownPlatforms[key];

if (!pkg) {
	console.error(
		`ktx is not available for platform ${process.platform}-${process.arch}\n` +
			`Supported platforms: linux (x64, arm64), darwin (x64, arm64), win32 (x64, arm64)`,
	);
	process.exit(1);
}

// Resolve the binary path relative to the main package
const binaryName = process.platform === "win32" ? "ktx.exe" : "ktx";
const binaryPath = path.join(__dirname, "..", pkg, binaryName);

const child = spawn(binaryPath, process.argv.slice(2), {
	stdio: "inherit",
	windowsHide: true,
});

child.on("exit", (code, signal) => {
	if (code !== null) process.exit(code);
	// On signal exit, forward the signal
	process.kill(process.pid, signal);
});

// Handle errors spawning the process
child.on("error", (err) => {
	console.error(`Failed to spawn ktx binary: ${err.message}`);
	console.error(`Expected at: ${binaryPath}`);
	console.error("Try reinstalling: npm install @kirimemail/ktx");
	process.exit(1);
});
