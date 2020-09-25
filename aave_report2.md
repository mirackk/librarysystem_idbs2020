## aave的lendingpool存入代币后的账户变化  



用户的balance数量，主要由对应的atokens contract中的balanceOf()函数给出  



![aave.vpd](C:\Users\87264\Downloads\aave.vpd.jpg)



流程主要分为两大类，当用户没有redirect过时，也就是自己的redirectionAddress指向自己，以及指向其他用户。

1.iRA==address(0)
		只需要简单的使用calculate()计算出本金会参与生成多少利息。
2.iRA!=address(0)
		此时会输入currentBalance+redirectBalance作为参数输入计算。往后的流程说明了redirectBalance是怎么出现的以及如何被改变的。redirectBalance会被一个update()改变，update()负责更新用户的redirectBalance，如果用户的interest没有被redirected，不会发生任何更新。是否知道有没有被redirected通过interestRedirectionAddresses[]的map记录，可以在外界调用函数更改它。
		与此同时，update()会被mintOnDeposit()函数调用，它会通过另一个合约lendingpool中的Deposit()函数，也就是存入时被调用。只要得到了redirectBalance就可以像1.一样使用calculate()来计算结果了。



calculateCumulatedBalanceInternal（x）在计算增值的时候公式如下:



<a href="https://www.codecogs.com/eqnedit.php?latex=\frac{\frac{\frac{RAY}{2}&plus;(x*RATIO*y)&plus;\frac{z}{2}}{z}&plus;\frac{RATIO}{2}}{RATIO}" target="_blank"><img src="https://latex.codecogs.com/gif.latex?\frac{\frac{\frac{RAY}{2}&plus;(x*RATIO*y)&plus;\frac{z}{2}}{z}&plus;\frac{RATIO}{2}}{RATIO}" title="\frac{\frac{\frac{RAY}{2}+(x*RATIO*y)+\frac{z}{2}}{z}+\frac{RATIO}{2}}{RATIO}" /></a>



两个常量 RAY 1e27; RATIO 1e9; 
三个变量 x为输入参数;
				y由lendingpoolcore中的对应代币（比如eth）的NormalizedIncome记录给出，数据存放在CoreLibrary.sol合约里。 
				z为一个map，userIndexes[user]，它赋值只有两种情况，当redirectedBalances==0，也就是用户未发生过redirect行为时，并且调用了对应atokens合约中的reset函数，z为0;否则，z=y，也就是同一个NormalizedIncome.
通常可以认为z==y，由此可以化简公式：



<a href="https://www.codecogs.com/eqnedit.php?latex=\frac{RAY}{2*RATIO*y}&plus;x&plus;0.5*RATIO&plus;0.5" target="_blank"><img src="https://latex.codecogs.com/gif.latex?\frac{RAY}{2*RATIO*y}&plus;x&plus;0.5*RATIO&plus;0.5" title="\frac{RAY}{2*RATIO*y}+x+0.5*RATIO+0.5" /></a>
