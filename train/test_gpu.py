
import os
#os.environ['CUDA_VISIBLE_DEVICES'] = "0"  # 使用第几个GPU， 0是第一个
import tensorflow as tf
import numpy as np
hello=tf.constant('hhh')
sess=tf.Session()
print (sess.run(hello))