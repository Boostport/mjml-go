import "fastestsmallesttextencoderdecoder-encodeinto/EncoderDecoderTogether.min.js";
import { compile } from "./lib";

import { setup, runnable } from "@suborbital/runnable";

const decoder = new TextDecoder();

function run_e(payload, ident) {
  // Imports will be injected by the runtime
  setup(this.imports, ident);

  const input = decoder.decode(payload);

  let encodedJSON;

  try {
    const decodedJSON = JSON.parse(input);
    const result = compile(decodedJSON);
    encodedJSON = JSON.stringify(result);
  } catch (err) {
    encodedJSON = JSON.stringify({
      error: {
        message: err.message,
      },
    });
  }

  runnable.returnResult(encodedJSON);
}

export { run_e };
