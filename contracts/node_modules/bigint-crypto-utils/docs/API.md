# bigint-crypto-utils - v3.2.2

## Table of contents

### Functions

- [abs](API.md#abs)
- [bitLength](API.md#bitlength)
- [eGcd](API.md#egcd)
- [gcd](API.md#gcd)
- [isProbablyPrime](API.md#isprobablyprime)
- [lcm](API.md#lcm)
- [max](API.md#max)
- [min](API.md#min)
- [modInv](API.md#modinv)
- [modPow](API.md#modpow)
- [prime](API.md#prime)
- [primeSync](API.md#primesync)
- [randBetween](API.md#randbetween)
- [randBits](API.md#randbits)
- [randBitsSync](API.md#randbitssync)
- [randBytes](API.md#randbytes)
- [randBytesSync](API.md#randbytessync)
- [toZn](API.md#tozn)

## Functions

### abs

▸ **abs**(`a`): `number` \| `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |

#### Returns

`number` \| `bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:1

___

### bitLength

▸ **bitLength**(`a`): `number`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |

#### Returns

`number`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:3

___

### eGcd

▸ **eGcd**(`a`, `b`): `Egcd`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `b` | `number` \| `bigint` |

#### Returns

`Egcd`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:10

___

### gcd

▸ **gcd**(`a`, `b`): `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `b` | `number` \| `bigint` |

#### Returns

`bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:12

___

### isProbablyPrime

▸ **isProbablyPrime**(`w`, `iterations?`, `disableWorkers?`): `Promise`<`boolean`\>

The test first tries if any of the first 250 small primes are a factor of the input number and then passes several
iterations of Miller-Rabin Probabilistic Primality Test (FIPS 186-4 C.3.1)

**`Throws`**

RangeError if w<0

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `w` | `number` \| `bigint` | `undefined` | A positive integer to be tested for primality |
| `iterations` | `number` | `16` | The number of iterations for the primality test. The value shall be consistent with Table C.1, C.2 or C.3 of FIPS 186-4 |
| `disableWorkers` | `boolean` | `false` | Disable the use of workers for the primality test |

#### Returns

`Promise`<`boolean`\>

A promise that resolves to a boolean that is either true (a probably prime number) or false (definitely composite)

#### Defined in

[src/ts/isProbablyPrime.ts:20](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/isProbablyPrime.ts#L20)

___

### lcm

▸ **lcm**(`a`, `b`): `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `b` | `number` \| `bigint` |

#### Returns

`bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:14

___

### max

▸ **max**(`a`, `b`): `number` \| `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `b` | `number` \| `bigint` |

#### Returns

`number` \| `bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:16

___

### min

▸ **min**(`a`, `b`): `number` \| `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `b` | `number` \| `bigint` |

#### Returns

`number` \| `bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:18

___

### modInv

▸ **modInv**(`a`, `n`): `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `n` | `number` \| `bigint` |

#### Returns

`bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:20

___

### modPow

▸ **modPow**(`b`, `e`, `n`): `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `b` | `number` \| `bigint` |
| `e` | `number` \| `bigint` |
| `n` | `number` \| `bigint` |

#### Returns

`bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:22

___

### prime

▸ **prime**(`bitLength`, `iterations?`): `Promise`<`bigint`\>

A probably-prime (Miller-Rabin), cryptographically-secure, random-number generator.
The browser version uses web workers to parallelise prime look up. Therefore, it does not lock the UI
main process, and it can be much faster (if several cores or cpu are available).
The node version can also use worker_threads if they are available (enabled by default with Node 11 and
and can be enabled at runtime executing node --experimental-worker with node >=10.5.0).

**`Throws`**

RangeError if bitLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `bitLength` | `number` | `undefined` | The required bit length for the generated prime |
| `iterations` | `number` | `16` | The number of iterations for the Miller-Rabin Probabilistic Primality Test |

#### Returns

`Promise`<`bigint`\>

A promise that resolves to a bigint probable prime of bitLength bits.

#### Defined in

[src/ts/prime.ts:28](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/prime.ts#L28)

___

### primeSync

▸ **primeSync**(`bitLength`, `iterations?`): `bigint`

A probably-prime (Miller-Rabin), cryptographically-secure, random-number generator.
The sync version is NOT RECOMMENDED since it won't use workers and thus it'll be slower and may freeze thw window in browser's javascript. Please consider using prime() instead.

**`Throws`**

RangeError if bitLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `bitLength` | `number` | `undefined` | The required bit length for the generated prime |
| `iterations` | `number` | `16` | The number of iterations for the Miller-Rabin Probabilistic Primality Test |

#### Returns

`bigint`

A bigint probable prime of bitLength bits.

#### Defined in

[src/ts/prime.ts:109](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/prime.ts#L109)

___

### randBetween

▸ **randBetween**(`max`, `min?`): `bigint`

Returns a cryptographically secure random integer between [min,max].

**`Throws`**

RangeError if max <= min

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `max` | `bigint` | Returned value will be <= max |
| `min` | `bigint` | Returned value will be >= min |

#### Returns

`bigint`

A cryptographically secure random bigint between [min,max]

#### Defined in

[src/ts/randBetween.ts:14](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/randBetween.ts#L14)

___

### randBits

▸ **randBits**(`bitLength`, `forceLength?`): `Promise`<`Uint8Array` \| `Buffer`\>

Secure random bits for both node and browsers. Node version uses crypto.randomFill() and browser one self.crypto.getRandomValues()

**`Throws`**

RangeError if bitLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `bitLength` | `number` | `undefined` | The desired number of random bits |
| `forceLength` | `boolean` | `false` | Set to true if you want to force the output to have a specific bit length. It basically forces the msb to be 1 |

#### Returns

`Promise`<`Uint8Array` \| `Buffer`\>

A Promise that resolves to a UInt8Array/Buffer (Browser/Node.js) filled with cryptographically secure random bits

#### Defined in

[src/ts/randBits.ts:13](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/randBits.ts#L13)

___

### randBitsSync

▸ **randBitsSync**(`bitLength`, `forceLength?`): `Uint8Array` \| `Buffer`

Secure random bits for both node and browsers. Node version uses crypto.randomFill() and browser one self.crypto.getRandomValues()

**`Throws`**

RangeError if bitLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `bitLength` | `number` | `undefined` | The desired number of random bits |
| `forceLength` | `boolean` | `false` | Set to true if you want to force the output to have a specific bit length. It basically forces the msb to be 1 |

#### Returns

`Uint8Array` \| `Buffer`

A Uint8Array/Buffer (Browser/Node.js) filled with cryptographically secure random bits

#### Defined in

[src/ts/randBits.ts:43](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/randBits.ts#L43)

___

### randBytes

▸ **randBytes**(`byteLength`, `forceLength?`): `Promise`<`Uint8Array` \| `Buffer`\>

Secure random bytes for both node and browsers. Node version uses crypto.randomBytes() and browser one self.crypto.getRandomValues()

**`Throws`**

RangeError if byteLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `byteLength` | `number` | `undefined` | The desired number of random bytes |
| `forceLength` | `boolean` | `false` | Set to true if you want to force the output to have a bit length of 8*byteLength. It basically forces the msb to be 1 |

#### Returns

`Promise`<`Uint8Array` \| `Buffer`\>

A promise that resolves to a UInt8Array/Buffer (Browser/Node.js) filled with cryptographically secure random bytes

#### Defined in

[src/ts/randBytes.ts:13](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/randBytes.ts#L13)

___

### randBytesSync

▸ **randBytesSync**(`byteLength`, `forceLength?`): `Uint8Array` \| `Buffer`

Secure random bytes for both node and browsers. Node version uses crypto.randomFill() and browser one self.crypto.getRandomValues()
This is the synchronous version, consider using the asynchronous one for improved efficiency.

**`Throws`**

RangeError if byteLength < 1

#### Parameters

| Name | Type | Default value | Description |
| :------ | :------ | :------ | :------ |
| `byteLength` | `number` | `undefined` | The desired number of random bytes |
| `forceLength` | `boolean` | `false` | Set to true if you want to force the output to have a bit length of 8*byteLength. It basically forces the msb to be 1 |

#### Returns

`Uint8Array` \| `Buffer`

A UInt8Array/Buffer (Browser/Node.js) filled with cryptographically secure random bytes

#### Defined in

[src/ts/randBytes.ts:54](https://github.com/juanelas/bigint-crypto-utils/blob/f4ea3d5/src/ts/randBytes.ts#L54)

___

### toZn

▸ **toZn**(`a`, `n`): `bigint`

#### Parameters

| Name | Type |
| :------ | :------ |
| `a` | `number` \| `bigint` |
| `n` | `number` \| `bigint` |

#### Returns

`bigint`

#### Defined in

node_modules/bigint-mod-arith/dist/index.d.ts:24
