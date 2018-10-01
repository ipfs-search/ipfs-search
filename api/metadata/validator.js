const isIPFS = require('is-ipfs')

function validate(cid) {
    // Validate cid, true for valid, false for invalid
    return isIPFS.cid(cid)
}

function make_validator() {
    // Create a validator
    return {
        Validate: validate
    }

}

module.exports = {
    Validator: make_validator
}

