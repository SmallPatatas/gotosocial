# 进阶设置

提供进阶设置选项是为了让管理员能够根据自己的喜好调整实例。

这些设置已设置为合理的默认值，所以大多数服务器管理员不需要更改或考虑它们。

**如果你不知道自己在做什么，修改这些设置可能会导致实例出错**。

## 设置

```yaml
#############################
##### 进阶设置 #####
#############################

# 与HTTP超时、安全性、Cookie等相关的进阶设置。
#
# 只有在你了解自己在做什么的情况下才调整这些设置！
#
# 大多数用户不需要（也不应该）修改这些设置，因为它们被设为合理的默认值，改变可能导致问题。
#
# 不过，这些设置提供给服务器管理员用于性能或安全原因的调整。

# 字符串。GoToSocial设置的Cookie的SameSite属性值。
# 默认设置为 'lax' 以确保OIDC流程不会中断，这通常是可以的。
# 如果你希望加强实例对抗CSRF攻击，并且不介意某些登录相关操作可能中断，可以将其设置为 'strict'。
#
# 关于此设置的概述，请参见：
# https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite
#
# 选项: ["lax", "strict"]
# 默认: "lax"
advanced-cookies-samesite: "lax"

# 整数。允许单个IP地址在5分钟内对每个路由分组的请求数量。
# 如果超出此数量，将返回429 HTTP错误代码。
#
# 如果你发现需要调整此限制，是因为它经常被超出，你应首先验证 `trusted-proxies` 配置是否正确。
# 在许多情况下，超出速率限制是因为你的实例将所有传入请求视为来自*相同的IP地址*（你可以通过查看实例日志中的客户端IP来验证）。
# 如果是这种情况，尝试在调整此速率限制设置*之前*将该IP地址添加到你的`trusted-proxies`中！
#
# 如果将此设置为0或更少，则完全禁用速率限制。
#
# 示例: [1000, 500, 0]
# 默认: 300
advanced-rate-limit-requests: 300

# 字符串数组。要从速率限制中排除的CIDR范围。
# CIDR范围内的任何IP的请求将不受速率限制，并且这些请求不会设置速率限制头。
#
# 对于IPv6，我们只考虑到/64的子网。如果你想开放更大的前缀，你需要列出多个前缀。
#
# 在以下示例情况下（可能还有很多其他情况），这可能很有用：
#
# 1. 你已设置使用API的自动化服务，但它频繁被限速，你信任它没有滥用实例。资源
#
# 2. 你和多人共用同一路由器/NAT登录同一实例，所以你们都有相同的IP地址，并且不断相互限速。
#
# 3. 你主要使用自己的家庭网络访问实例，并希望豁免家庭网络的速率限制。
#
# 调整此设置时需要小心，因为如果设置范围过宽，可能会使速率限制变得无用。如果不确定，建议宁少勿多，并根据需要调整。
#
# 示例: ["192.168.0.0/16", "2001:DB8:FACE:CAFE::/64"]
# 默认: []
advanced-rate-limit-exceptions: []

# 整数。每个CPU、每个路由分组允许的开放请求数量，以应用HTTP请求限制。
# 超出计算限制的请求将被保留在一个等待队列中，最长30秒然后处理或超时。
# 不在等待队列中的请求将返回状态503，并设置“Retry-After”头为30秒。
#
# 开放请求限制为可用CPU * 乘数；等待队列限制为限制 * 乘数。
#
# 乘数为8的示例值：
#
# 1 cpu = 08 开放, 064 等待
# 2 cpu = 16 开放, 128 等待
# 4 cpu = 32 开放, 256 等待
#
# 乘数为4的示例值：
#
# 1 cpu = 04 开放, 016 等待
# 2 cpu = 08 开放, 032 等待
# 4 cpu = 16 开放, 064 等待
#
# 乘数为8是合理的默认值，但对于运行在性能非常高的硬件上的实例，你可能希望增加它；对于使用非常慢的CPU的实例，你可能希望减少它。
#
# 如果将此设置为0或更少，将完全禁用HTTP请求限制。
#
# 示例: [8, 4, 9, 0]
# 默认: 8
advanced-throttling-multiplier: 8

# 持续时间。用于响应限速请求的“retry-after”头值的时间段。
# 最小分辨率为1秒。
#
# 示例: [30s, 10s, 5s, 1m]
# 默认: "30s"
advanced-throttling-retry-after: "30s"

# 整数。用于通过ActivityPub发送消息的固定协程数量的CPU倍数。
# 消息将被批量处理并推送到单一队列，倍数 * CPU数的协程将提取对垒中的消息并尝试发送。
# 这可以用于限制对外站收件箱的并发发布，防止当有很多关注者的账户发布贴文时实例CPU使用率激增。
#
# 如果将此设置为0或更少，无论CPU数量如何，都只会使用1个发送者。这可能在你有非常严格的网络或CPU限制时有用。
#
# 乘数为2的示例值（默认）：
#
# 1 cpu = 2 个并发发送者
# 2 cpu = 4 个并发发送者
# 4 cpu = 8 个并发发送者
#
# 乘数为4的示例值：
#
# 1 cpu = 4 个并发发送者
# 2 cpu = 8 个并发发送者
# 4 cpu = 16 个并发发送者
#
# 乘数<1的示例值：
#
# 1 cpu = 1 个并发发送者
# 2 cpu = 1 个并发发送者
# 4 cpu = 1 个并发发送者
advanced-sender-multiplier: 2

# 字符串数组。为实例设置Content-Security-Policy头时，要添加到'img-src'和'media-src'中的额外URI。
#
# 这可以用于在浏览器中查看实例页面和个人资料时，允许加载来自额外来源（如S3桶等）的资源。
#
# 由于非代理的S3存储将在实例启动时被探测以生成正确的Content-Security-Policy，你可能永远都不需要修改此设置，但把它包括在内是因为“可配置项（通常）越多越好”。
# 
# 参见: https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
#
# 示例: ["s3.example.org", "some-bucket-name.s3.example.org"]
# 默认: []
advanced-csp-extra-uris: []

# 字符串。用于此实例的HTTP请求头过滤模式。
#
# "block" -- 只有明确被请求头过滤规则阻止的请求会被拒绝（除非它们被明确允许）。
#
# "allow" -- 只有明确被请求头过滤规则允许的请求会被接受（除非它们被明确阻止）。
#             此模式被视为实验性功能，并且几乎肯定会破坏对你实例的访问，除非非常小心。
#
#   ""    -- 请求头过滤禁用。
#
# 有关阻止和允许模式的更多详细信息，请查看文档：
# https://docs.gotosocial.org/zh-cn/latest/admin/request_filtering_modes
#
# 选项: ["block", "allow", ""]
# 默认: ""
advanced-header-filter-mode: ""
```