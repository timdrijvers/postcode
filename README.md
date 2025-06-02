# Postcode 
Minimal parser to extract postal codes, house number and street names from BAG extract.

This is a very rough first implementation, offering basic functionality on how to parse the 'lvbag' extract 
and extract postal code, house number and street name. It then generates a compact hash-map stored in a gob file.

The parser does not check for validity of entries, such as the life time timestamps.
