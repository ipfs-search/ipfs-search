const fs = require('fs');

const rawdata = fs.readFileSync('mapping.json');
const mapping = JSON.parse(rawdata);

// console.log(mapping.ipfs_files_v7.mappings.properties);

function doIt(pre, props) {
  let result = [];

  Object.keys(props).forEach((key) => {
    const prop = props[key];
    let fullKey;

    if (pre) {
      fullKey = [pre, key].join('.');
    } else {
      fullKey = key;
    }

    if ('index' in prop && prop.index === false) {
      return;
    }

    if ('properties' in prop) {
      // Recurse
      result = result.concat(doIt(fullKey, prop.properties));
    } else {
      // Add
      result.push(fullKey);
    }
  });

  // console.log('result', pre, result);
  return result;
}

console.log(doIt('', mapping.ipfs_files_v7.mappings.properties));
