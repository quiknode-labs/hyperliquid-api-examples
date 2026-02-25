import { rpc, signHash, getMid } from "./client";

const COIN = "BTC";

async function main() {
  const mid = await getMid(COIN);
  if (mid === 0) {
    console.error(`Could not fetch ${COIN} mid price`);
    process.exit(1);
  }

  const sz = (11.0 / mid).toFixed(5);
  const buyPx = Math.floor(mid * 1.03).toString();

  console.log(`${COIN} mid: $${mid.toLocaleString()}`);
  console.log(`BUY ${sz} @ ${buyPx} (IOC, ~$${(parseFloat(sz) * mid).toFixed(2)} notional)`);

  const action = {
    type: "order",
    orders: [{ a: COIN, b: true, p: buyPx, s: sz, r: false, t: { limit: { tif: "Ioc" } } }],
    grouping: "na",
  };

  const res = await rpc("hl_buildOrder", { action });
  const sig = await signHash(res.result.hash);

  const result = await rpc("hl_sendOrder", {
    action: res.result.action || action,
    nonce: res.result.nonce,
    signature: sig,
  });

  console.log(JSON.stringify(result.result.exchangeResponse, null, 2));
}

main();
