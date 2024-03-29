# BruceCore Development Log

### 2019-08-20

关于颜色差异，今天查到其实有比欧式距离更官方的做法，关键字 CIE76、CIE94。

今天尝试了一下图片优化，几个问题：

1. 人眼能分辨出好坏，但算法似乎很难分辨出好坏。  
现在能很精确的统计出图片差异，而且统计了方差，理论上来说，越接近源图的，效果越好。  
事实上，如果是一种压缩算法，基本上可以按这个来区分，但如果是2种算法，譬如png 和 webp（类jpg），png是增加颜色抖动，减少颜色数量（256色）。
2. 同算法下，不同图片，其实也很难计算出压缩到哪一步是性价比最合适的，容量和效果之间缺乏量化的转换关系，我能很容易知道降低质量，带来的好处在降低，但不知道如何取舍。

然后就是，https://pngquant.org 似乎还可以。

### 2019-08-19

最近在处理页面分析，遇到个问题，就是做某些分析时，其实往往是一个复杂分类模型，就是可能有A、B、C3种分类方法，可以分成 A > B > C 这样的，其实有时也想看看 B > A > C，甚至有时就只想看 A ，这种动态分类。

昨天还想到几个图片的量化指标：

- 像素质量，就是 图片大小 / 像素大小，这样可以算出单像素质量，用于不同大小不同格式之间图片横向比对
- 颜色距离图，颜色距离就是HSV的欧氏距离，我们压缩了一张图片，和原图片对比，能得出一张每个像素距离源图的距离图出来，对这个图做量化分析，是否能知道压缩过了呢？
- 如果一张图片需要压缩为png-8，我们能不能算出颜色数量，通过一套算法，评估是否适合压缩，png-8就是256色，如果图片数量太多，或距离过大，其实是明显不适合抖动到256色的