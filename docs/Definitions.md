# Definitions

This document aims to create a shared understanding of what some terms mean in the context of the R3 documentation.

**Share (Content)** - The collection of data (link, text, image, or file) that the user is trying to share or receive using R3.

**Share ID** - The identifier that is used to both identify and encrypt a share; this is the part of the link after the hash (i.e.: `Il2.exa.mpl.ejR` is the Share ID in `https://rkt.one/#Ils.3j9.F2k.fKL`.)

**Share Hash** - The SHA-512 hash derived from the Share ID that is sent to our servers and stored in our DB.

**Share Key** - A SHA-256 encryption key derived using the Share ID and the Share Salt and used to generate the Encrypted Content from a Share.

**Encrypted (Share) Content** - The Share after it's encrypted using AES-GCM, a key derived from the Share ID, and a randomly generated IV .

**(Share) Key Salt** - A JavaScript Uint8Array(12) filled with random values that is used as salt during the Key derivation process.

**(Share) Hash Salt** - A JavaScript Uint8Array filled with hard-coded (!) random values that is used as salt during the Hash derivation process.

**(Share) Encryption IV** - A randomly generated IV used for the encryption.

**Expiration Time** - The expiration time in seconds before a share is deleted by a server. The countdown starts when the share is added to the DB.

**Expiration Views / Max Views** - The number of views or queries for a given share before the server removes it.
