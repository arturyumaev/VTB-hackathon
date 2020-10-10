const express = require('express')
const app = express()
const port = 8888

app.get('/api', (request, response) => {
    console.log(request);
    return response.send('Hello from Express!');
})

app.listen(port, (err) => {
    if (err) {
        return console.log('something bad happened', err)
    }
    console.log(`server is listening on ${port}`)
})
