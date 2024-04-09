# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [eupalinos.proto](#eupalinos-proto)
    - [ChannelInfo](#proto-ChannelInfo)
    - [Diexodos](#proto-Diexodos)
    - [Empty](#proto-Empty)
    - [EnqueueResponse](#proto-EnqueueResponse)
    - [Epistello](#proto-Epistello)
    - [InternalEpistello](#proto-InternalEpistello)
    - [MessageUpdate](#proto-MessageUpdate)
    - [QueueLength](#proto-QueueLength)
  
    - [Operation](#proto-Operation)
  
    - [Eupalinos](#proto-Eupalinos)
  
- [Scalar Value Types](#scalar-value-types)



<a name="eupalinos-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## eupalinos.proto



<a name="proto-ChannelInfo"></a>

### ChannelInfo
Message for specifying the channel name


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="proto-Diexodos"></a>

### Diexodos



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="proto-Empty"></a>

### Empty







<a name="proto-EnqueueResponse"></a>

### EnqueueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="proto-Epistello"></a>

### Epistello
Public Epistello message without traceid


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| data | [string](#string) |  |  |
| channel | [string](#string) |  |  |






<a name="proto-InternalEpistello"></a>

### InternalEpistello
Internal Epistello message with traceid


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| data | [string](#string) |  |  |
| channel | [string](#string) |  |  |
| traceid | [string](#string) |  |  |






<a name="proto-MessageUpdate"></a>

### MessageUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [Operation](#proto-Operation) |  |  |
| message | [InternalEpistello](#proto-InternalEpistello) |  |  |






<a name="proto-QueueLength"></a>

### QueueLength
Response message for getting the length of the queue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| length | [int32](#int32) |  |  |





 


<a name="proto-Operation"></a>

### Operation


| Name | Number | Description |
| ---- | ------ | ----------- |
| ENQUEUE | 0 |  |
| DEQUEUE | 1 |  |


 

 


<a name="proto-Eupalinos"></a>

### Eupalinos


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StreamQueueUpdates | [MessageUpdate](#proto-MessageUpdate) stream | [InternalEpistello](#proto-InternalEpistello) stream | Bidirectional Streaming for task updates between Eupalinos pods |
| EnqueueMessage | [Epistello](#proto-Epistello) | [EnqueueResponse](#proto-EnqueueResponse) | Unary RPC for epistello enqueueing |
| DequeueMessage | [ChannelInfo](#proto-ChannelInfo) | [Epistello](#proto-Epistello) | Unary RPC for epistello dequeueing |
| GetQueueLength | [ChannelInfo](#proto-ChannelInfo) | [QueueLength](#proto-QueueLength) | Unary RPC for getting the length of the queue |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

