const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("ProofOfScience", function() {
  let ProofOfScience;
  let proofOfScience;
  let addr1;
  let addrs;

  beforeEach(async function () {
    ProofOfScience = await ethers.getContractFactory("ProofOfScience");
    [addr1, ...addrs] = await ethers.getSigners();
    proofOfScience = await ProofOfScience.deploy();
    await proofOfScience.deployed();
  });

  describe("Deployment", function() {
    it("Initial token ID should be zero", async function() {
      expect(await proofOfScience.tokenID()).to.equal(0);
    });
  });

  describe("Transactions", function() {
    it("Should mint a token", async function() {
      await proofOfScience.connect(addr1).mint(addr1.address, "QmHash");
      expect(await proofOfScience.balanceOf(addr1.address, 0)).to.equal(1);
    });

    it("Should set token URI correctly", async function() {
      await proofOfScience.connect(addr1).mint(addr1.address, "QmHash");
      expect(await proofOfScience.uri(0)).to.equal("ipfs://QmHash");
    });

    it("Should increment token ID correctly", async function() {
      await proofOfScience.connect(addr1).mint(addr1.address, "QmHash1");
      expect(await proofOfScience.tokenID()).to.equal(1);
      await proofOfScience.connect(addr1).mint(addr1.address, "QmHash2");
      expect(await proofOfScience.tokenID()).to.equal(2);
    });

    it("Should mint a token if called by any address", async function() {
      const [addr2] = addrs;
      await proofOfScience.connect(addr2).mint(addr2.address, "QmHash");
      expect(await proofOfScience.balanceOf(addr2.address, 0)).to.equal(1);
    });
  });

  describe("CantBeEvil license", function() {
    it("Should return correct license URI", async function() {
      expect(await proofOfScience.getLicenseURI()).to.equal("ar://zmc1WTspIhFyVY82bwfAIcIExLFH5lUcHHUN0wXg4W8/1");
    });

    it("Should return correct license name", async function() {
      expect(await proofOfScience.getLicenseName()).to.equal("EXCLUSIVE");
    });
  });
});
