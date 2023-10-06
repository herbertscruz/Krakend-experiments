let random = 1;

module.exports = (req, res, next) => {
    if (random === 0) {
        random = 1;
        res.setHeader('Content-Type', 'application/json');
        res.status(500).send({ error: "Sample error message" });
    } else {
        random = 0;
        next();
    }

}