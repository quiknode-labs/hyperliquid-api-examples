"""Check builder fee approval status."""

import json
from client import rpc, wallet

res = rpc("hl_getMaxBuilderFee", {"user": wallet.address})
print(json.dumps(res["result"], indent=2))
