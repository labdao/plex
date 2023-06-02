// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import {LicenseVersion, CantBeEvil} from "@a16z/contracts/licenses/CantBeEvil.sol";

/**
 * @title ProofOfScience
 * @dev This contract mints ERC-1155 tokens with IPFS URIs. Any account is allowed to mint new tokens.
 * The contract adheres to the license specified in the CantBeEvil contract, which is CBE EXCLUSIVE in this case.
 */
contract ProofOfScience is ERC1155, CantBeEvil {
    uint256 public tokenID = 0;
    mapping (uint256 => string) private _tokenURIs;

    // Base URI for all tokens
    string private _baseURI = "ipfs://";

    /**
     * @dev Contract constructor that sets the base URI for all tokens in the contract.
     */
    constructor() ERC1155(_baseURI) CantBeEvil(LicenseVersion.EXCLUSIVE) {}

    /**
     * @dev Mints a new token and assigns it to `account`, 
     * increasing the total supply.
     *
     * @param account Recipient of the token minting.
     * @param tokenURI The IPFS URI of the associated token data.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     * - The token doesn't exist, `id` must not exist in other token's URI.
     */
    function mint(address account, string memory tokenURI) public {
        _mint(account, tokenID, 1, "");
        _setTokenURI(tokenID, tokenURI);
        tokenID = tokenID + 1;
    }

    /**
     * @dev Returns the IPFS URI for a given token ID
     *
     * @param _id uint256 ID of the token to query
     * @return URI string
     */
    function uri(uint256 _id) public view override returns (string memory) {
        return string(abi.encodePacked(_baseURI, _tokenURIs[_id]));
    }

    /**
     * @dev Internal function to set the token URI for a given token
     *
     * @param tokenId uint256 ID of the token to set its URI
     * @param _tokenURI string URI to assign
     */
    function _setTokenURI(uint256 tokenId, string memory _tokenURI) internal virtual {
        _tokenURIs[tokenId] = _tokenURI;
    }

    /**
     * @dev Override supportsInterface to use both ERC1155 and CantBeEvil's implementations.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC1155, CantBeEvil) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}