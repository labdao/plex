#!/usr/bin/env node
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const semver_1 = __importDefault(require("semver"));
const chalk_1 = __importDefault(require("chalk"));
const SUPPORTED_NODE_VERSIONS = ["^14.0.0", "^16.0.0", "^18.0.0"];
if (!semver_1.default.satisfies(process.version, SUPPORTED_NODE_VERSIONS.join(" || "))) {
    console.warn(chalk_1.default.yellow.bold(`WARNING:`), `You are using a version of Node.js that is not supported, and it may work incorrectly, or not work at all. See https://hardhat.org/nodejs-versions`);
    console.log();
    console.log();
}
require("./cli");
//# sourceMappingURL=bootstrap.js.map