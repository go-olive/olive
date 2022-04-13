# go-olive
save lives



uploader 

* youtube 
    * daily limitation
    * Concurrent number
    * tryout count
* bilibili



all the uploader has succeeded 

using hook func to 

* move to new folder
* or delete it

and delete it from savers



use cron to delete new folder's content to save space



recorder作为生产者生成upload task

upload center作为管理者管理所有task

* task 内容
    * 上传的文件名（唯一索引）
    * 上传的平台
    * 重试的次数
    * ...

根据task的上传平台进行任务的分发

分发到具体每个上传平台的处理中心

* 每个上传平台有自己的消费者
* 消费者数量可配置
* 上传结束后要有回调函数

upload center根据回调情况做一个处理

* 成功 - 是否全部都成功

    * 都成功执行成功回调
        * 移动位置
        * 删除文件

* 失败 - 是否重传

    * 满足重传条件重传

    * 若发生特定错误

        * api调用超过限额，停止相应平台消费者，

            并过段时间唤醒

        * auth认证失败，直接停用

平台消费者要支持指令退出
