function validTimestamp() {
    return Math.round(new Date().getTime()/1000) + 600;
}

function expiredTimestamp() {
    return Math.round(new Date().getTime()/1000) - 600;
}

module.exports = {
    validTimestamp,
    expiredTimestamp,
};