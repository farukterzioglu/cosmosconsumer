`~/.gaiaconsumer`  
`~/.gaiaconsumer/config.json`  
```json
{
    "accountlist" : [
        "cosomos1", 
        "cosomos2"
    ]
}
```
`~/.gaiaconsumer/leveldb`  
`./gaiaconsumer --start-height 5 --port 8080`  

Flags;   
```bash
--data-dir (~/.gaiaconsumer)
--port (8080)
--start-height
```

```bash
curl localhost:8080/block/5
curl localhost:8080/tx/txhash123
```


https://docs.tendermint.com/master/rpc/#/Websocket/subscribe  


```json
{ "jsonrpc": "2.0", "method": "subscribe", "params": ["tm.event='Tx' And transfer.sender='cosmos17yj45mrgvwezj9jlnhcgd2wr33ufqlcflj0xxv'"], "id": 1 }
{ "jsonrpc": "2.0", "method": "unsubscribe", "params": ["tm.event='Tx' And transfer.sender='cosmos17yj45mrgvwezj9jlnhcgd2wr33ufqlcflj0xxv'"], "id": 1 }

{ "jsonrpc": "2.0", "method": "subscribe", "params": ["tm.event='Tx' And transfer.recipient='cosmos1vm975pe9hrghrdfg29ssdnkh0pl4m4vnckcaw0'"], "id": 1 }
{ "jsonrpc": "2.0", "method": "unsubscribe", "params": ["tm.event='Tx' And transfer.recipient='cosmos1vm975pe9hrghrdfg29ssdnkh0pl4m4vnckcaw0'"], "id": 1 }

{ "jsonrpc": "2.0", "method": "subscribe", "params": ["tm.event='Tx'"], "id": 1 }
{ "jsonrpc": "2.0", "method": "unsubscribe", "params": ["tm.event='Tx'"], "id": 1 }

{ "jsonrpc": "2.0", "method": "subscribe", "params": ["tm.event='NewBlock' and tx.height > 0"], "id": 1 }
{ "jsonrpc": "2.0", "method": "unsubscribe", "params": ["tm.event='NewBlock' and tx.height > 0"], "id": 1 }
```