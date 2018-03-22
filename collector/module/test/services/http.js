const http = require('http');


const port = process.argv[2] || 8000;

console.log('listening to port',port);


http.createServer( async (request, response) => {

    response.writeHead(200, {'Content-type':'text/plan'});
    response.write(`${port}`);
    response.end( );
}).listen(port);

