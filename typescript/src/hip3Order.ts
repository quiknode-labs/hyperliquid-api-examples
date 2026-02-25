import { rpc, signHash } from "./client";

const COIN = "xyz:SILVER";
const DEX = COIN.split(":")[0];

async function main() {
  const midRes = await fetch("https://api.hyperliquid.xyz/info", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ type: "allMids", dex: DEX }),
  });
  const mids: any = await midRes.json();
  let mid = parseFloat(mids[COIN] || "0");
  if (mid === 0) {
    console.log(`Could not fetch ${COIN} mid price, using fallback`);
    mid = 78.0;
  }

  const sz = (10.0 / mid).toFixed(2);
  const buyPx = (mid * 1.03).toFixed(2);

  console.log(`${COIN} mid: $${mid.toFixed(2)}`);
  console.log(`BUY ${sz} @ ${buyPx} (IOC, ~$${(parseFloat(sz) * mid).toFixed(2)} notional)\n`);

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
