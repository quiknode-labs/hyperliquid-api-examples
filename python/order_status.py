"""Check what happened to an order by OID."""

import json
import sys
from client import rpc, wallet

oid = int(sys.argv[1]) if len(sys.argv) > 1 else 0
if oid == 0:
    print("Usage: python3 order_status.py <oid>")
    print("Get OIDs from: python3 open_orders.py")
    sys.exit(1)

res = rpc("hl_orderStatus", {"user": wallet.address, "oid": oid})
result = res["result"]

if result["status"] == "unknownOid":
    print(f"Order {oid}: not found")
    print(f"  {result['explanation']}")
    sys.exit(0)

spot = " [SPOT]" if result.get("isSpot") else ""
print(f"Order {oid} on {result['name']}{spot}: {result['status']}")
print(f"  {result['explanation']}")
if result.get("order"):
    print(f"  Details: {json.dumps(result['order'], indent=2)}")
