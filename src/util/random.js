const possibleCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
const charactersLength = possibleCharacters.length;

module.exports = {
    generateRandomString: function(length) {
        let result = "";
        for (let i = 0; i < length; i++) {
            result += possibleCharacters.charAt(Math.floor(Math.random() * charactersLength));
        }
        return result;
    }
};