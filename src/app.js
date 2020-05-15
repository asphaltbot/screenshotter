let fastify = undefined;
const route = require('./routes/takeScreenshot');
const config = require("../config.json");

async function initialise() {
   let fastifyConfig = {
      logger: true
   };

   if (config.UseSSL) {
      fastifyConfig.key = config.SSLConfig.KeyPath;
      fastifyConfig.cert = config.SSLConfig.CertPath;
   }

   fastify = require('fastify')(fastifyConfig);

   fastify.post("/screenshot", route.takeScreenshot);

   fastify.listen(config.Port, config.Host, (err, address) => {
      if (err) throw err;
      fastify.log.info(`API listening on ${address}`);
   });
}

(async () => {
   await route.init();
   await initialise();
})()