const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("ProofOfScience", function() {
  let ProofOfScience;
  let proofOfScience;
  let owner;
  let addr1;
  let addr2;
  let addrs;

  beforeEach(async function () {
    ProofOfScience = await ethers.getContractFactory("ProofOfScience");
    [owner, addr1, addr2, ...addrs] = await ethers.getSigners();
    proofOfScience = await ProofOfScience.deploy();
    await proofOfScience.deployed();
  });

  describe("Deployment", function() {
    it("Should set the right owner", async function() {
      expect(await proofOfScience.owner()).to.equal(owner.address);
    });

    it("Initial token ID should be zero", async function() {
      expect(await proofOfScience.tokenID()).to.equal(0);
    });
  });

  describe("Transactions", function() {
    it("Should mint a token", async function() {
      await proofOfScience.connect(owner).mint(addr1.address, "QmHash");
      expect(await proofOfScience.balanceOf(addr1.address, 0)).to.equal(1);
    });

    it("Should set token URI correctly", async function() {
      await proofOfScience.connect(owner).mint(addr1.address, "QmHash");
      expect(await proofOfScience.uri(0)).to.equal("ipfs://QmHash");
    });

    it("Should increment token ID correctly", async function() {
      await proofOfScience.connect(owner).mint(addr1.address, "QmHash1");
      expect(await proofOfScience.tokenID()).to.equal(1);
      await proofOfScience.connect(owner).mint(addr1.address, "QmHash2");
      expect(await proofOfScience.tokenID()).to.equal(2);
    });

    it("Should not mint a token if called by non-owner", async function() {
      await expect(proofOfScience.connect(addr1).mint(addr1.address, "QmHash")).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should transfer contract ownership correctly", async function() {
      await proofOfScience.connect(owner).transferContractOwnership(addr1.address);
      expect(await proofOfScience.owner()).to.equal(addr1.address);
    });

    it("Should fail to transfer ownership if not called by owner", async function() {
      await expect(proofOfScience.connect(addr2).transferContractOwnership(addr1.address)).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });
});