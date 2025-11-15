// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title SimpleCounter
 * @dev 一个极简的计数器智能合约，仅包含核心功能
 */
contract SimpleCounter {
    uint256 public count;
    
    // 事件：当计数器增加时触发
    event CountIncremented(uint256 newCount);
    
    // 构造函数：初始化计数器为0
    constructor() {
        count = 0;
    }
    
    /**
     * @dev 增加计数器值（每次增加1）
     */
    function increment() public {
        count += 1;
        emit CountIncremented(count);
    }
    
    /**
     * @dev 获取当前计数器值
     * @return 当前计数值
     */
    function getCount() public view returns (uint256) {
        return count;
    }
}