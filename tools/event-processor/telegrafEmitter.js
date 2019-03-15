const net = require('net');
const fs = require('fs');

module.exports = class TelegrafEmitter {
    constructor(socketPath) {
        this.socketPath = socketPath;
        if(!fs.existsSync(this.socketPath)) {
            throw Error("Socket doesn't exist! Start telegraf first.")
        }
        this.socket = new net.Socket();
        this.socket.connect(this.socketPath, () => {
            console.log("Socket connected");
        });
    }

    emit(companyId, packet) {
        if (!packet.sourceId || !packet.events) {
            return;
        }

        const events = packet.events.map(e => {
            e.sourceId = packet.sourceId;
            e.companyId = companyId;
            return e;
        });

        for(let event of events) {
            console.log(event);
            this.socket.write(JSON.stringify(event));
            this.socket.write('\r\n');
        }
    }
}