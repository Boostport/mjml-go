import * as http from "http";
import { compile } from "./lib";

http
  .createServer(async (request, response) => {
    let result;

    if (request.method !== "POST") {
      result = {
        error: {
          message: "Only POST requests are accepted",
        },
      };
    } else {
      const buffers = [];
      for await (const chunk of request) {
        buffers.push(chunk);
      }
      const data = Buffer.concat(buffers).toString();

      try {
        const input = JSON.parse(data);
        result = compile(input);
      } catch (e) {
        result = {
          error: {
            message: e.message,
          },
        };
      }
    }

    const encoded = JSON.stringify(result);

    response.writeHead(200, { "Content-Type": "application/json" });
    response.end(encoded, "utf-8");
  })
  .listen(8888);
