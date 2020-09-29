pragma solidity ^0.6.0;
pragma experimental ABIEncoderV2;

contract Proxy {
    constructor() public{
    }
    
    event CallResult(uint256 index, bool success);

    function call(address[] memory addrs,bytes[] memory parameters) public {
        for (uint256 i = 0; i < addrs.length; i++){
            (bool success, ) = addrs[i].call(parameters[i]);
            emit CallResult(i,success);
        }
    }
}
