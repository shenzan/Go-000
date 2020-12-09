### 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么
* 应该Warp，dao层的sql.ErrNoRows不应该传递到service层，service和dao需要解耦
* 另外sql.ErrNoRows在service层看来不一定是异常，所以需要传递一定的信息给service，我这里做了个新的接口
* 其实生产中使用github.com/pkg/errors可能是更好的选择，但是作为作业，还是自己实现一下比较好