# API doc

This is API documentation for Topic append API. This is generated by `httpdoc`. Don't edit by hand.

## Table of contents

- [[200] POST /meta/channel](#200-post-metachannel)


## [200] POST /meta/channel

Register new topic

### Request



Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Content-Type | application/json |  |





Request example

```

{
	"channel": "govim"
}

```


### Response

Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Content-Type | application/json |  |





Response example

```
{"channel":"govim","successful":true}
```


