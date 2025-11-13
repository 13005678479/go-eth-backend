// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title SimpleCounter
 * @dev 一个简单的计数器智能合约，用于演示以太坊智能合约开发
 */
contract SimpleCounter {
    uint256 private count;
    address public owner;
    
    // 事件定义
    event CountIncremented(uint256 newCount, address indexed by);
    event CountDecremented(uint256 newCount, address indexed by);
    event CountReset(uint256 newCount, address indexed by);
    
    // 修饰器：只有合约所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
    
    // 构造函数
    constructor() {
        owner = msg.sender;
        count = 0;
    }
    
    /**
     * @dev 增加计数器值
     * @param _amount 要增加的数量
     */
    function increment(uint256 _amount) public {
        count += _amount;
        emit CountIncremented(count, msg.sender);
    }
    
    /**
     * @dev 减少计数器值
     * @param _amount 要减少的数量
     */
    function decrement(uint256 _amount) public {
        require(count >= _amount, "Cannot decrement below zero");
        count -= _amount;
        emit CountDecremented(count, msg.sender);
    }
    
    /**
     * @dev 重置计数器为0
     */
    function reset() public onlyOwner {
        count = 0;
        emit CountReset(count, msg.sender);
    }
    
    /**
     * @dev 获取当前计数器值
     * @return 当前计数值
     */
    function getCount() public view returns (uint256) {
        return count;
    }
    
    /**
     * @dev 获取合约所有者地址
     * @return 所有者地址
     */
    function getOwner() public view returns (address) {
        return owner;
    }
    
    /**
     * @dev 转移合约所有权
     * @param _newOwner 新的所有者地址
     */
    function transferOwnership(address _newOwner) public onlyOwner {
        require(_newOwner != address(0), "New owner cannot be zero address");
        owner = _newOwner;
    }
    
    /**
     * @dev 自毁合约（仅限所有者调用）
     */
    function destroy() public onlyOwner {
        selfdestruct(payable(owner));
    }
}