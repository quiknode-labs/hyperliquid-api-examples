import { exchange, signHash, address } from "./client";

const COIN = "HYPE";

async function main() {
  console.log(`Closing ${COIN} position for ${address}\n`);

  // Optional: custom slippage (default 3%, range 0.1%-10%)
  // const res = await exchange({ action: {...}, slippage: 0.05 });  // 5% slippage

  const res = await exchange({
    action: {
      type: "closePosition",
      asset: COIN,
      user: address,
    },
  });

  const ctx = res.closePositionContext || {};
  console.log(`Position: ${ctx.positionSize} ${ctx.positionSide}`);
  console.log(`Close: ${ctx.closeSide} ${ctx.closeSize} @ ${ctx.slippedPrice}`);

  const sig = await signHash(res.hash);

  const result = await exchange({
    action: res.action,
    nonce: res.nonce,
    signature: sig,
  });

  console.log(JSON.stringify(result.exchangeResponse, null, 2));
  console.log("\nPosition closed.");
}

main();
