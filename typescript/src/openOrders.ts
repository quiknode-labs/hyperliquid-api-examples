import { rpc, wallet } from "./client";

async function main() {
  const res = await rpc("hl_openOrders", { user: wallet.address });
  const result = res.result;

  console.log(`Open orders: ${result.count}`);
  for (const order of result.orders) {
    const side = order.side === "B" ? "BUY" : "SELL";
    const spot = order.isSpot ? " [SPOT]" : "";
    console.log(
      `  ${order.name}${spot} ${side} ${order.sz} @ ${order.limitPx} (OID: ${order.oid})`
    );
  }

  if (result.count > 0) {
    console.log("\nCancel actions by asset:");
    for (const [name, action] of Object.entries<any>(result.cancelActions.byAsset)) {
      console.log(`  ${name}: ${action.cancels.length} order(s) — pass to hl_buildCancel`);
    }

    console.log("\nTo cancel ALL orders, pass cancelActions.all to hl_buildCancel:");
    console.log(JSON.stringify(result.cancelActions.all, null, 2));
  }
}

main();
