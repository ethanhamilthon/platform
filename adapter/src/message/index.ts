import { connect, NatsConnection, StringCodec, Codec, NatsError } from "nats";

export default class NatsMessage {
  private _nc: NatsConnection | null = null;
  private _decoder: Codec<string>;

  constructor() {
    this._decoder = StringCodec();
  }

  // connect creates new connection to nats server
  connect = async (addr: string) => {
    try {
      this._nc = await connect({ servers: addr });
      console.log("Connected to NATS");
    } catch (error) {
      console.error("Failed to connect to NATS:", error);
      throw error;
    }
  };

  // addsub adds new subscriber to the topic, callback can return data to response or null
  addSub(
    topic: string,
    callback: (err: NatsError | null, msg: string) => Promise<string | void>
  ) {
    if (!this._nc) {
      throw new Error("Not connected to NATS");
    }

    this._nc.subscribe(topic, {
      callback: (err, msg) => {
        callback(err, msg.json()).then((value) => {
          if (value === undefined) return;
          msg.respond(value);
        });
      },
    });
  }

  // publish pushs new message to the topic
  publish(topic: string, data: any) {
    if (!this._nc) {
      throw new Error("Not connected to NATS");
    }

    const encodedData = this._decoder.encode(JSON.stringify(data));
    this._nc.publish(topic, encodedData);
  }

  close = async (): Promise<void> => {
    if (this._nc) {
      await this._nc.close();
      console.log("NATS connection closed");
    }
  };

  request = async (topic: string, data: any) => {
    if (!this._nc) {
      throw new Error("Not connected to NATS");
    }
    const reply = await this._nc.request(topic, JSON.stringify(data));
    return this._decoder.decode(reply.data);
  };

  errorHandler(cb: (message: string) => Promise<string | void>) {
    return async (
      err: NatsError | null,
      msg: string
    ): Promise<string | void> => {
      if (err !== null) {
        this.publish("adapter:error", err);
      } else {
        return cb(msg);
      }
    };
  }
}
