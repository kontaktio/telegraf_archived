const net = require('net');
const fs = require('fs');

module.exports = class TelegrafEmitter {
    constructor(socketPath) {
        this.socketPath = socketPath;
        if(!fs.existsSync(this.socketPath)) {
            throw Error("Socket doesn't exist! Start telegraf first.")
        }
        this.socket = new net.Socket();
        this.connect();

        this.processors = {
            3: this.processV3,
            4: this.processV4
        }
    }

    connect() {
        this.socket.connect(this.socketPath, () => {
            console.log("Socket connected");
        });
    }

    emit(companyId, packet) {
        if (!packet.sourceId || !packet.events) {
            return;
        }

        const events = this.processors[packet.version](companyId, packet);

        for(let event of events) {
            try {
                this.socket.write(JSON.stringify(event));
                this.socket.write('\r\n');
            } catch (e) {
                console.error(`Exception while sending data: ${e.message}. Quitting`);
                process.exit();
            }
        }
    }

    processV3(companyId, packet) {
        return packet.events.map(e => {
            e.sourceId = packet.sourceId;
            e.companyId = companyId;
            return e;
        });
    }

    processV4(companyId, packet) {
        return packet
            .events
            .filter(e => e.ble)
            .map(e => {
                return {
                    ...e.ble,
                    companyId: companyId,
                    sourceId: e.sourceId,
                    rssi: e.rssi,
                    timestamp: e.timestamp
                }
            })
    }
};