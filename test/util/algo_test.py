from sysl.util import algo
from sysl.util import debug

import unittest

class TestUtil(unittest.TestCase):
	
	def test_topo_sort(self):
	  G = {
	    1: {2, 3},
	    2: {4},
	    3: {4},
	    4: {6},
	    5: {3, 6},
	  }
	  S = [n for (_, n) in sorted(algo.topo_sort(G))]
	  
	  assert S == [1, 5, 2, 3, 4, 6], S
 
if __name__ == '__main__':
	unittest.main()