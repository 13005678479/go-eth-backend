const solc = require('solc');
const fs = require('fs');
const path = require('path');

// 读取合约文件
const contractPath = path.join(__dirname, '../contracts/Counter.sol');
const source = fs.readFileSync(contractPath, 'utf8');

// 编译配置
const input = {
  language: 'Solidity',
  sources: {
    'Counter.sol': {
      content: source
    }
  },
  settings: {
    outputSelection: {
      '*': {
        '*': ['abi', 'evm.bytecode']
      }
    }
  }
};

// 编译合约
const output = JSON.parse(solc.compile(JSON.stringify(input)));

// 检查编译错误
if (output.errors) {
  console.error('Compilation errors:', output.errors);
  process.exit(1);
}

// 提取合约信息
const contractName = 'SimpleCounter';
const contractData = output.contracts['Counter.sol'][contractName];

if (!contractData) {
  console.error(`Contract ${contractName} not found in output`);
  process.exit(1);
}

// 创建 build 目录
const buildDir = path.join(__dirname, '../contracts/build');
if (!fs.existsSync(buildDir)) {
  fs.mkdirSync(buildDir, { recursive: true });
}

// 保存 ABI
fs.writeFileSync(
  path.join(buildDir, 'Counter.abi'),
  JSON.stringify(contractData.abi, null, 2)
);

// 保存字节码
fs.writeFileSync(
  path.join(buildDir, 'Counter.bin'),
  contractData.evm.bytecode.object
);

console.log('✅ Contract compiled successfully!');
console.log('📁 ABI saved to: contracts/build/Counter.abi');
console.log('📁 Bytecode saved to: contracts/build/Counter.bin');