const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("ProofOfScience", (m) => {
  const proofOfScience = m.contract("ProofOfScience", []);
  return { proofOfScience };
});