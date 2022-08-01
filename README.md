
magicbox


## 错误处理

遵循https://lailin.xyz/post/go-training-03.html

1. 如果是调用本项目其他函数出现的问题，应该直接返回err，如果想携带额外信息就使用errors.WithMessage
2. 如果是应用程序本身逻辑预期之内出现的问题，一般使用 errors.New 或 errors.Errorf 返回错误
3. 如果调用其他库（标准库、企业公共库、开源第三方库等）的函数发生了错误，应用程序应该使用 errors.Wrap获取堆栈信息。只需要在错误第一次出现时使用，且在编写基础库和被大量引用的第三方库中不使用，避免堆栈信息重复。
4. 使用 errors.Is判断错误类型，使用errors.As做赋值
