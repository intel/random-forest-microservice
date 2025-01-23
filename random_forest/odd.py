# Copyright (c) 2025 Intel Corporation.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

import numpy as np 

def automated_optimal_binning(data_list):
    """
    Data mining algorithm implementation
    """

    class HistData:
        """
        Class to hold intermediate data values
        """

        def __init__(
            self, cost=0, bin_count=0, bin_width=0, bin_array=0, bin_data=0
        ):
            self.cost = cost
            self.bin_count = bin_count
            self.bin_width = bin_width
            self.bin_array = bin_array
            self.bin_data = bin_data
    try:
        if data_list is None:
            print("Invalid input data")
            return (None, None, None, None)
        # optimization: calculate max/min outside of loop
        data_orig_max = max(data_list)
        data_orig_min = min(data_list)

        # optimization: calculate mean/stdev outside of loop
        stat_mean = np.mean(data_list)
        stat_stdev = np.std(data_list)
        
        #print(np.mean(data_list), np.std(data_list))
    
        mul_factor = 1
        t = mul_factor * stat_stdev
    
        # create new data list within 1 stddev
        data_new=data_list[abs(data_list - stat_mean) <= t]
        
        data_new_max = max(data_new)
        data_new_min = min(data_new)

        step = 0.005
        start_value = 0.01
        end_value = 0.1
        cost_map = {}

        for resval in np.arange(start_value, end_value, step):
            bin_width = (data_new_max - data_new_min) * resval
            xbins = np.arange(data_orig_min, data_orig_max, bin_width)
            hist_data, bin_edges = np.histogram(data_list, bins=xbins)
            # using ediff1d to calculate diff of consective numbers in an array
            cost = max(abs(np.ediff1d(hist_data))) / bin_width
            cost_map[bin_width] = HistData(
                cost, len(bin_edges), bin_width, xbins, hist_data
            )
        
        # Find out hist distribution with lowest cost
        width_opt = min(cost_map.keys(), key=(lambda k: cost_map[k].cost))

        return (
            cost_map[width_opt].bin_count,
            cost_map[width_opt].bin_width,
            cost_map[width_opt].bin_array,
            cost_map[width_opt].bin_data,
        )
    except Exception as e:
        print(e)