import { exchange, signHash } from "./client";

const COIN = "BTC";
const SIZE = "0.00011";

async function main() {
  console.log(`Market BUY ${SIZE} ${COIN}\n`);

  // Optional: custom slippage (default 3%, range 0.1%-10%)
  // const res = await exchange({ action: {...}, slippage: 0.05 });  // 5% slippage

  const res = await exchange({
    action: {
      type: "order",
      orders: [{ asset: COIN, side: "buy", size: SIZE, tif: "market" }],
    },
  });

  const computedPrice = res.action.orders[0].p;
  console.log(`Computed price (mid + slippage, default 3%): ${computedPrice}`);
  console.log(`Builder fee: ${res.builderFee}`);

  const sig = await signHash(res.hash);

  const result = await exchange({
    action: res.action,
    nonce: res.nonce,
    signature: sig,
  });

  console.log(JSON.stringify(result.exchangeResponse, null, 2));
}

main();
