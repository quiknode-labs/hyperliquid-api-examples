"""View open orders with enriched asset info and pre-built cancel actions."""

import json
from client import rpc, wallet

res = rpc("hl_openOrders", {"user": wallet.address})
result = res["result"]

print(f"Open orders: {result['count']}")
for order in result["orders"]:
    side = "BUY" if order["side"] == "B" else "SELL"
    spot = " [SPOT]" if order.get("isSpot") else ""
    print(f"  {order['name']}{spot} {side} {order['sz']} @ {order['limitPx']} (OID: {order['oid']})")

if result["count"] > 0:
    print(f"\nCancel actions by asset:")
    for name, action in result["cancelActions"]["byAsset"].items():
        count = len(action["cancels"])
        print(f"  {name}: {count} order(s) — pass to hl_buildCancel")

    print(f"\nTo cancel ALL orders, pass cancelActions.all to hl_buildCancel:")
    print(json.dumps(result["cancelActions"]["all"], indent=2))
