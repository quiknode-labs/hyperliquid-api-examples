import { rpc, signHash, getMid } from "./client";

const COIN = "BTC";

async function main() {
  const mid = await getMid(COIN);
  if (mid === 0) {
    console.error(`Could not fetch ${COIN} mid price`);
    process.exit(1);
  }

  const sz = (11.0 / mid).toFixed(5);
  const restPx = Math.floor(mid * 0.97).toString();

  console.log(`${COIN} mid: $${mid.toLocaleString()}`);
  console.log(`Placing resting BUY ${sz} @ ${restPx} (GTC, 3% below mid)\n`);

  const orderAction = {
    type: "order",
    orders: [{ a: COIN, b: true, p: restPx, s: sz, r: false, t: { limit: { tif: "Gtc" } } }],
    grouping: "na",
  };

  let res = await rpc("hl_buildOrder", { action: orderAction });
  let sig = await signHash(res.result.hash);

  const result = await rpc("hl_sendOrder", {
    action: res.result.action || orderAction,
    nonce: res.result.nonce,
    signature: sig,
  });

  const exchange = result.result.exchangeResponse;
  const statuses = exchange?.response?.data?.statuses || [];

  let oid: number | null = null;
  for (const s of statuses) {
    if (s.resting) {
      oid = s.resting.oid;
      break;
    }
  }

  if (oid === null) {
    console.error("Could not extract OID from resting order");
    console.log(JSON.stringify(exchange, null, 2));
    process.exit(1);
  }

  console.log(`Order resting (OID: ${oid})`);
  console.log("Cancelling...\n");

  const cancelAction = {
    type: "cancel",
    cancels: [{ a: COIN, o: oid }],
  };

  res = await rpc("hl_buildCancel", { action: cancelAction });
  sig = await signHash(res.result.hash);

  const cancelResult = await rpc("hl_sendCancel", {
    action: cancelAction,
    nonce: res.result.nonce,
    signature: sig,
  });

  console.log(JSON.stringify(cancelResult.result.exchangeResponse, null, 2));
  console.log("\nOrder cancelled.");
}

main();
