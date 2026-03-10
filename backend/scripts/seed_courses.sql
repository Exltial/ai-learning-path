-- ============================================================
-- AI 学习平台 - 课程种子数据
-- 创建 5 个示例课程，包含完整章节
-- ============================================================

-- 课程 1: Python 编程入门 (讲师：prof_chen)
INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, estimated_hours, price, is_published, enrollment_count, rating, created_at) VALUES
('30000000-0000-0000-0000-000000000001', 
 'Python 编程入门', 
 '从零开始学习 Python 编程语言，掌握基础语法、数据结构、函数和面向对象编程。适合编程零基础的学习者。',
 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=400',
 '10000000-0000-0000-0000-000000000001',
 'Programming',
 'beginner',
 40,
 0.00,
 true,
 156,
 4.8,
 NOW() - INTERVAL '60 days');

-- 课程 1 的章节
INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at) VALUES
('40000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '欢迎与课程介绍', '了解课程结构和学习目标', '欢迎来到 Python 编程入门课程！本课程将带你从零开始学习 Python。', 'https://video.example.com/python-intro.mp4', 180, 1, true, NOW() - INTERVAL '60 days'),
('40000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'Python 环境搭建', '安装 Python 和配置开发环境', '学习如何安装 Python 解释器，配置 IDE，创建第一个 Hello World 程序。', 'https://video.example.com/python-setup.mp4', 600, 2, true, NOW() - INTERVAL '59 days'),
('40000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000001', '变量与数据类型', '理解 Python 中的变量和基本数据类型', '深入学习 Python 的变量命名规则、整数、浮点数、字符串、布尔值等基本数据类型。', 'https://video.example.com/python-variables.mp4', 900, 3, false, NOW() - INTERVAL '58 days'),
('40000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000001', '条件语句', '掌握 if-else 条件判断', '学习如何使用 if、elif、else 语句进行条件判断，理解布尔逻辑。', 'https://video.example.com/python-conditionals.mp4', 720, 4, false, NOW() - INTERVAL '57 days'),
('40000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000001', '循环结构', '学习 for 和 while 循环', '掌握 for 循环遍历序列，while 循环条件执行，以及 break 和 continue 的使用。', 'https://video.example.com/python-loops.mp4', 840, 5, false, NOW() - INTERVAL '56 days'),
('40000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000001', '函数定义', '创建可复用的代码块', '学习如何定义函数、传递参数、返回值，理解作用域概念。', 'https://video.example.com/python-functions.mp4', 960, 6, false, NOW() - INTERVAL '55 days'),
('40000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000001', '列表与元组', 'Python 序列类型详解', '深入学习列表和元组的创建、访问、修改和常用方法。', 'https://video.example.com/python-lists.mp4', 780, 7, false, NOW() - INTERVAL '54 days'),
('40000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000001', '字典与集合', '键值对和无序集合', '掌握字典的键值对操作和集合的去重、交集、并集等运算。', 'https://video.example.com/python-dicts.mp4', 720, 8, false, NOW() - INTERVAL '53 days'),
('40000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000001', '面向对象编程', '类和对象的基础', '理解类、对象、属性、方法的概念，学习封装、继承、多态。', 'https://video.example.com/python-oop.mp4', 1200, 9, false, NOW() - INTERVAL '52 days'),
('40000000-0000-0000-0000-00000000000a', '30000000-0000-0000-0000-000000000001', '文件操作与异常处理', '读写文件和错误处理', '学习文件的打开、读写、关闭，以及 try-except 异常处理机制。', 'https://video.example.com/python-files.mp4', 660, 10, false, NOW() - INTERVAL '51 days');

-- 课程 2: 机器学习基础 (讲师：prof_chen)
INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, estimated_hours, price, is_published, enrollment_count, rating, created_at) VALUES
('30000000-0000-0000-0000-000000000002', 
 '机器学习基础', 
 '系统学习机器学习的核心概念和算法，包括监督学习、无监督学习和模型评估方法。',
 'https://images.unsplash.com/photo-1555949963-ff9fe0c870eb?w=400',
 '10000000-0000-0000-0000-000000000001',
 'Artificial Intelligence',
 'intermediate',
 60,
 99.00,
 true,
 89,
 4.9,
 NOW() - INTERVAL '55 days');

-- 课程 2 的章节
INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at) VALUES
('40000000-0000-0000-0000-00000000000b', '30000000-0000-0000-0000-000000000002', '机器学习概述', '什么是机器学习及其应用场景', '介绍机器学习的基本概念、发展历程和主要应用领域。', 'https://video.example.com/ml-intro.mp4', 540, 1, true, NOW() - INTERVAL '55 days'),
('40000000-0000-0000-0000-00000000000c', '30000000-0000-0000-0000-000000000002', 'Python 数据科学库', 'NumPy、Pandas、Matplotlib 基础', '学习数据科学必备的 Python 库，掌握数组操作、数据处理和可视化。', 'https://video.example.com/ml-libraries.mp4', 1080, 2, false, NOW() - INTERVAL '54 days'),
('40000000-0000-0000-0000-00000000000d', '30000000-0000-0000-0000-000000000002', '数据预处理', '数据清洗和特征工程', '掌握缺失值处理、异常值检测、特征缩放、编码等预处理技术。', 'https://video.example.com/ml-preprocessing.mp4', 900, 3, false, NOW() - INTERVAL '53 days'),
('40000000-0000-0000-0000-00000000000e', '30000000-0000-0000-0000-000000000002', '线性回归', '预测连续值的经典算法', '深入理解线性回归原理，学习模型训练、评估和优化方法。', 'https://video.example.com/ml-linear-regression.mp4', 960, 4, false, NOW() - INTERVAL '52 days'),
('40000000-0000-0000-0000-00000000000f', '30000000-0000-0000-0000-000000000002', '逻辑回归', '分类问题的基础算法', '学习逻辑回归的数学原理和在二分类、多分类问题中的应用。', 'https://video.example.com/ml-logistic-regression.mp4', 840, 5, false, NOW() - INTERVAL '51 days'),
('40000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000002', '决策树与随机森林', '树模型集成方法', '理解决策树的分裂策略，掌握随机森林的集成学习思想。', 'https://video.example.com/ml-decision-trees.mp4', 1020, 6, false, NOW() - INTERVAL '50 days'),
('40000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000002', '支持向量机', '最大间隔分类器', '学习 SVM 的原理、核函数技巧和在复杂分类问题中的应用。', 'https://video.example.com/ml-svm.mp4', 900, 7, false, NOW() - INTERVAL '49 days'),
('40000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000002', 'K-Means 聚类', '无监督学习的入门算法', '掌握 K-Means 聚类算法原理、实现和实际应用案例。', 'https://video.example.com/ml-kmeans.mp4', 720, 8, false, NOW() - INTERVAL '48 days'),
('40000000-0000-0000-0000-000000000013', '30000000-0000-0000-0000-000000000002', '模型评估与调优', '交叉验证和超参数优化', '学习准确率、精确率、召回率等评估指标，掌握网格搜索调参方法。', 'https://video.example.com/ml-evaluation.mp4', 780, 9, false, NOW() - INTERVAL '47 days'),
('40000000-0000-0000-0000-000000000014', '30000000-0000-0000-0000-000000000002', '实战项目：房价预测', '综合应用所学知识', '通过完整的房价预测项目，整合课程所学知识解决实际问题。', 'https://video.example.com/ml-project.mp4', 1800, 10, false, NOW() - INTERVAL '46 days');

-- 课程 3: Web 开发实战 (讲师：prof_liu)
INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, estimated_hours, price, is_published, enrollment_count, rating, created_at) VALUES
('30000000-0000-0000-0000-000000000003', 
 'Web 开发实战', 
 '从零构建现代 Web 应用，学习 HTML、CSS、JavaScript 和后端 API 开发。',
 'https://images.unsplash.com/photo-1547658719-da2b51169166?w=400',
 '10000000-0000-0000-0000-000000000002',
 'Web Development',
 'beginner',
 50,
 0.00,
 true,
 203,
 4.7,
 NOW() - INTERVAL '50 days');

-- 课程 3 的章节
INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at) VALUES
('40000000-0000-0000-0000-000000000015', '30000000-0000-0000-0000-000000000003', 'Web 开发基础', '了解 Web 工作原理', '学习 HTTP 协议、客户端 - 服务器模型、浏览器工作原理。', 'https://video.example.com/web-intro.mp4', 420, 1, true, NOW() - INTERVAL '50 days'),
('40000000-0000-0000-0000-000000000016', '30000000-0000-0000-0000-000000000003', 'HTML5 核心', '构建网页结构', '深入学习 HTML5 语义化标签、表单、多媒体元素。', 'https://video.example.com/web-html.mp4', 780, 2, true, NOW() - INTERVAL '49 days'),
('40000000-0000-0000-0000-000000000017', '30000000-0000-0000-0000-000000000003', 'CSS3 样式', '美化网页外观', '掌握 CSS 选择器、盒模型、Flexbox、Grid 布局、动画效果。', 'https://video.example.com/web-css.mp4', 960, 3, false, NOW() - INTERVAL '48 days'),
('40000000-0000-0000-0000-000000000018', '30000000-0000-0000-0000-000000000003', 'JavaScript 基础', '让网页动起来', '学习变量、函数、DOM 操作、事件处理等 JavaScript 核心概念。', 'https://video.example.com/web-js.mp4', 1080, 4, false, NOW() - INTERVAL '47 days'),
('40000000-0000-0000-0000-000000000019', '30000000-0000-0000-0000-000000000003', '异步编程', 'Ajax 和 Fetch API', '理解异步编程概念，掌握 Promise、async/await、API 调用。', 'https://video.example.com/web-async.mp4', 840, 5, false, NOW() - INTERVAL '46 days'),
('40000000-0000-0000-0000-00000000001a', '30000000-0000-0000-0000-000000000003', '后端基础', 'Node.js 入门', '学习 Node.js 环境、Express 框架、路由和中间件。', 'https://video.example.com/web-nodejs.mp4', 900, 6, false, NOW() - INTERVAL '45 days'),
('40000000-0000-0000-0000-00000000001b', '30000000-0000-0000-0000-000000000003', '数据库集成', 'MongoDB 基础', '掌握 NoSQL 数据库概念，学习 MongoDB CRUD 操作。', 'https://video.example.com/web-mongodb.mp4', 720, 7, false, NOW() - INTERVAL '44 days'),
('40000000-0000-0000-0000-00000000001c', '30000000-0000-0000-0000-000000000003', '用户认证', 'JWT 和会话管理', '实现用户注册、登录、JWT Token 认证和权限控制。', 'https://video.example.com/web-auth.mp4', 660, 8, false, NOW() - INTERVAL '43 days'),
('40000000-0000-0000-0000-00000000001d', '30000000-0000-0000-0000-000000000003', '部署上线', '将应用发布到互联网', '学习 Git 版本控制、Docker 容器化、云服务器部署。', 'https://video.example.com/web-deploy.mp4', 600, 9, false, NOW() - INTERVAL '42 days'),
('40000000-0000-0000-0000-00000000001e', '30000000-0000-0000-0000-000000000003', '综合项目：博客系统', '完整的全栈项目开发', '从零开始构建一个功能完整的博客系统，整合所有所学知识。', 'https://video.example.com/web-project.mp4', 2400, 10, false, NOW() - INTERVAL '41 days');

-- 课程 4: 数据结构与算法 (讲师：prof_huang)
INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, estimated_hours, price, is_published, enrollment_count, rating, created_at) VALUES
('30000000-0000-0000-0000-000000000004', 
 '数据结构与算法', 
 '深入理解核心数据结构和经典算法，提升编程能力和面试技巧。',
 'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=400',
 '10000000-0000-0000-0000-000000000003',
 'Computer Science',
 'intermediate',
 45,
 129.00,
 true,
 127,
 4.9,
 NOW() - INTERVAL '45 days');

-- 课程 4 的章节
INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at) VALUES
('40000000-0000-0000-0000-00000000001f', '30000000-0000-0000-0000-000000000004', '算法复杂度分析', '时间复杂度和空间复杂度', '学习大 O 表示法，分析算法效率，理解最好/最坏/平均情况。', 'https://video.example.com/dsa-complexity.mp4', 540, 1, true, NOW() - INTERVAL '45 days'),
('40000000-0000-0000-0000-000000000020', '30000000-0000-0000-0000-000000000004', '数组与字符串', '基础线性结构', '掌握数组的操作、字符串处理技巧和常用算法。', 'https://video.example.com/dsa-arrays.mp4', 720, 2, false, NOW() - INTERVAL '44 days'),
('40000000-0000-0000-0000-000000000021', '30000000-0000-0000-0000-000000000004', '链表', '动态线性结构', '理解单链表、双链表、循环链表的实现和应用。', 'https://video.example.com/dsa-linkedlist.mp4', 840, 3, false, NOW() - INTERVAL '43 days'),
('40000000-0000-0000-0000-000000000022', '30000000-0000-0000-0000-000000000004', '栈与队列', '特殊线性表', '学习栈和队列的特性、实现和应用场景。', 'https://video.example.com/dsa-stack-queue.mp4', 660, 4, false, NOW() - INTERVAL '42 days'),
('40000000-0000-0000-0000-000000000023', '30000000-0000-0000-0000-000000000004', '哈希表', '高效查找结构', '理解哈希函数、冲突解决和哈希表的应用。', 'https://video.example.com/dsa-hash.mp4', 780, 5, false, NOW() - INTERVAL '41 days'),
('40000000-0000-0000-0000-000000000024', '30000000-0000-0000-0000-000000000004', '树与二叉树', '层次结构数据', '掌握树的遍历、二叉搜索树、平衡树的概念。', 'https://video.example.com/dsa-trees.mp4', 960, 6, false, NOW() - INTERVAL '40 days'),
('40000000-0000-0000-0000-000000000025', '30000000-0000-0000-0000-000000000004', '堆与优先队列', '特殊树结构', '学习最小堆、最大堆和优先队列的实现与应用。', 'https://video.example.com/dsa-heap.mp4', 720, 7, false, NOW() - INTERVAL '39 days'),
('40000000-0000-0000-0000-000000000026', '30000000-0000-0000-0000-000000000004', '图论基础', '网络结构表示', '理解图的表示方法、遍历算法和最短路径。', 'https://video.example.com/dsa-graphs.mp4', 1020, 8, false, NOW() - INTERVAL '38 days'),
('40000000-0000-0000-0000-000000000027', '30000000-0000-0000-0000-000000000004', '排序算法', '经典排序方法', '深入学习冒泡、选择、插入、快速、归并、堆排序等算法。', 'https://video.example.com/dsa-sorting.mp4', 1080, 9, false, NOW() - INTERVAL '37 days'),
('40000000-0000-0000-0000-000000000028', '30000000-0000-0000-0000-000000000004', '搜索与动态规划', '高级算法技巧', '掌握 DFS、BFS、回溯法和动态规划的核心思想。', 'https://video.example.com/dsa-dp.mp4', 1200, 10, false, NOW() - INTERVAL '36 days');

-- 课程 5: 深度学习入门 (讲师：prof_huang)
INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, estimated_hours, price, is_published, enrollment_count, rating, created_at) VALUES
('30000000-0000-0000-0000-000000000005', 
 '深度学习入门', 
 '学习神经网络和深度学习基础，掌握 TensorFlow 和 PyTorch 框架。',
 'https://images.unsplash.com/photo-1555255707-c07966088b7b?w=400',
 '10000000-0000-0000-0000-000000000003',
 'Artificial Intelligence',
 'advanced',
 70,
 199.00,
 true,
 64,
 4.8,
 NOW() - INTERVAL '40 days');

-- 课程 5 的章节
INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at) VALUES
('40000000-0000-0000-0000-000000000029', '30000000-0000-0000-0000-000000000005', '深度学习概述', '神经网络的发展与应用', '了解深度学习的历史、神经网络基本原理和主要应用领域。', 'https://video.example.com/dl-intro.mp4', 600, 1, true, NOW() - INTERVAL '40 days'),
('40000000-0000-0000-0000-00000000002a', '30000000-0000-0000-0000-000000000005', '神经网络基础', '感知机与多层网络', '学习神经元模型、激活函数、前向传播和反向传播算法。', 'https://video.example.com/dl-nn.mp4', 1080, 2, false, NOW() - INTERVAL '39 days'),
('40000000-0000-0000-0000-00000000002b', '30000000-0000-0000-0000-000000000005', 'TensorFlow 基础', 'Google 深度学习框架', '掌握 TensorFlow 2.x 的基本用法和张量操作。', 'https://video.example.com/dl-tensorflow.mp4', 900, 3, false, NOW() - INTERVAL '38 days'),
('40000000-0000-0000-0000-00000000002c', '30000000-0000-0000-0000-000000000005', 'PyTorch 基础', 'Facebook 深度学习框架', '学习 PyTorch 的动态图机制和模型构建方法。', 'https://video.example.com/dl-pytorch.mp4', 900, 4, false, NOW() - INTERVAL '37 days'),
('40000000-0000-0000-0000-00000000002d', '30000000-0000-0000-0000-000000000005', '卷积神经网络', '图像处理的核心技术', '深入理解 CNN 的卷积层、池化层和经典网络架构。', 'https://video.example.com/dl-cnn.mp4', 1200, 5, false, NOW() - INTERVAL '36 days'),
('40000000-0000-0000-0000-00000000002e', '30000000-0000-0000-0000-000000000005', '循环神经网络', '序列数据处理', '学习 RNN、LSTM、GRU 的原理和应用。', 'https://video.example.com/dl-rnn.mp4', 1080, 6, false, NOW() - INTERVAL '35 days'),
('40000000-0000-0000-0000-00000000002f', '30000000-0000-0000-0000-000000000005', '迁移学习', '利用预训练模型', '掌握迁移学习思想和 Fine-tuning 技巧。', 'https://video.example.com/dl-transfer.mp4', 780, 7, false, NOW() - INTERVAL '34 days'),
('40000000-0000-0000-0000-000000000030', '30000000-0000-0000-0000-000000000005', '生成对抗网络', '创意 AI 的基础', '理解 GAN 的生成器和判别器对抗训练机制。', 'https://video.example.com/dl-gan.mp4', 960, 8, false, NOW() - INTERVAL '33 days'),
('40000000-0000-0000-0000-000000000031', '30000000-0000-0000-0000-000000000005', '模型优化与部署', '提升性能和上线', '学习模型压缩、量化、蒸馏和部署方法。', 'https://video.example.com/dl-deploy.mp4', 840, 9, false, NOW() - INTERVAL '32 days'),
('40000000-0000-0000-0000-000000000032', '30000000-0000-0000-0000-000000000005', '实战项目：图像分类', '完整深度学习项目', '使用 CNN 实现图像分类系统，整合课程所学知识。', 'https://video.example.com/dl-project.mp4', 2100, 10, false, NOW() - INTERVAL '31 days');
