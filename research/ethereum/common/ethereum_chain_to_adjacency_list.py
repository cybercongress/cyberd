import sys

from web3.auto import w3

"""
Collects ethereum transactions as adjacency list.

Output - txt file with next structure:
FROM TO VALUE
str str int
...
0xF84e2C88af13BB807a5FAc27c936a3ee6FfCb6E4 0xEB6D43Fe241fb2320b5A3c9BE9CDfD4dd8226451 1150660000000000000
0x07ac11Be43bFE9223D74FE97AE5e54116Ba70d52 0xEB6D43Fe241fb2320b5A3c9BE9CDfD4dd8226451 240077040000000000
0xaec588A285AB9A84E23ffF171510FF19ae8faB1F 0xEB6D43Fe241fb2320b5A3c9BE9CDfD4dd8226451 6994572000000000000
0x8e1af970F090778CF7c56Ac956a5A2762A9f9fAc 0x9A755332D874c893111207b0b220Ce2615cd036F 1995000000000000000
0x6483598aFDfA001eC71234B91da2ac9284e2f048 0x9A755332D874c893111207b0b220Ce2615cd036F 1000500000000000000
0x69b148395Ce0015C13e36BFfBAd63f49EF874E03 0x6748F50f686bfbcA6Fe8ad62b22228b87F31ff2b 0


Usage:
1) Start local Parity full node
2) Run Script with two params: start block  and end blocknumbers ex:
python ethereum_chain_to_adjacency_list.py 6235000 6235009
"""

first_block_to_download = int(sys.argv[1])
last_block_to_download = int(sys.argv[2])

file_name = "../data/{}-{}_blocks_data".format(first_block_to_download, last_block_to_download)

print("")
print("-----------------------------------------------")
print("About to download {}-{} blocks".format(first_block_to_download, last_block_to_download))
print("")
result_data_file = open(file_name, "w")

for block_number in range(first_block_to_download, last_block_to_download + 1):

    print("Downloading {} block".format(block_number))
    block = w3.eth.getBlock(block_number, True)

    for tx in block.transactions:
        if tx.to is not None:
            result_data_file.write("{} {} {}\r\n".format(tx.to, getattr(tx, 'from'), tx.value))
        else:
            result_data_file.write("{} {} {}\r\n".format(tx.creates, getattr(tx, 'from'), tx.value))

result_data_file.close()
print("Finished to download data into {}".format(file_name))
print("-----------------------------------------------")
print("")
