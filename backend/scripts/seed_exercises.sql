-- ============================================================
-- AI 学习平台 - 练习题种子数据
-- 每个课程 10+ 道题，5 种题型：multiple_choice, coding, fill_blank, true_false, essay
-- ============================================================

-- ============ 课程 1: Python 编程入门 练习题 ============

-- 第 1 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', '课程目标理解', '以下哪项不是本课程的学习目标？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "掌握 Python 基础语法", "is_correct": false}, {"text": "学会开发操作系统", "is_correct": true}, {"text": "理解面向对象编程", "is_correct": false}, {"text": "能够编写实用程序", "is_correct": false}]', '{"correct_option": 1}', 1, NOW() - INTERVAL '60 days'),
('50000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000001', 'Python 应用领域', 'Python 可以应用于哪些领域？（多选）', 'multiple_choice', 'easy', 15, 3, 90, '[{"text": "Web 开发", "is_correct": true}, {"text": "数据科学", "is_correct": true}, {"text": "人工智能", "is_correct": true}, {"text": "硬件制造", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 2, NOW() - INTERVAL '60 days');

-- 第 2 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000002', 'Python 版本检查', '在命令行中检查 Python 版本的命令是？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "python --version", "is_correct": true}, {"text": "python -v", "is_correct": false}, {"text": "python version", "is_correct": false}, {"text": "check python", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '59 days'),
('50000000-0000-0000-0000-000000000004', '40000000-0000-0000-0000-000000000002', '第一个程序', '编写一个 Python 程序，输出 "Hello, World!"', 'coding', 'easy', 20, 5, 300, 'null', '{"expected_output": "Hello, World!", "test_cases": [{"input": "", "expected": "Hello, World!"}]}', 2, NOW() - INTERVAL '59 days');

-- 第 3 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000005', '40000000-0000-0000-0000-000000000003', '变量命名规则', '以下哪个变量名是合法的？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "2name", "is_correct": false}, {"text": "_name", "is_correct": true}, {"text": "my-name", "is_correct": false}, {"text": "class", "is_correct": false}]', '{"correct_option": 1}', 1, NOW() - INTERVAL '58 days'),
('50000000-0000-0000-0000-000000000006', '40000000-0000-0000-0000-000000000003', '数据类型判断', 'type(3.14) 的返回值是？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "float", "alternatives": ["<class ''float''>"]}', 2, NOW() - INTERVAL '58 days'),
('50000000-0000-0000-0000-000000000007', '40000000-0000-0000-0000-000000000003', '字符串操作', 'Python 中字符串是不可变类型。', 'true_false', 'easy', 5, 2, 30, 'null', '{"answer": true}', 3, NOW() - INTERVAL '58 days');

-- 第 4 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000008', '40000000-0000-0000-0000-000000000004', '条件语句语法', 'Python 中使用什么关键字表示"否则"？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "else", "alternatives": ["ELSE", "Else"]}', 1, NOW() - INTERVAL '57 days'),
('50000000-0000-0000-0000-000000000009', '40000000-0000-0000-0000-000000000004', '比较运算符', '编写代码判断一个数是否为偶数', 'coding', 'medium', 20, 5, 300, 'null', '{"expected_function": "is_even", "test_cases": [{"input": "4", "expected": "True"}, {"input": "7", "expected": "False"}]}', 2, NOW() - INTERVAL '57 days'),
('50000000-0000-0000-0000-00000000000a', '40000000-0000-0000-0000-000000000004', '逻辑运算符', 'and 运算符的优先级高于 or 运算符。', 'true_false', 'easy', 5, 2, 30, 'null', '{"answer": true}', 3, NOW() - INTERVAL '57 days');

-- 第 5 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000000b', '40000000-0000-0000-0000-000000000005', 'for 循环语法', '以下哪个是正确的 for 循环语法？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "for i in range(5):", "is_correct": true}, {"text": "for (i=0; i<5; i++)", "is_correct": false}, {"text": "foreach i in 5:", "is_correct": false}, {"text": "loop i from 1 to 5:", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '56 days'),
('50000000-0000-0000-0000-00000000000c', '40000000-0000-0000-0000-000000000005', '计算累加和', '使用循环计算 1 到 100 的和', 'coding', 'medium', 25, 5, 400, 'null', '{"expected_output": "5050", "test_cases": [{"input": "", "expected": "5050"}]}', 2, NOW() - INTERVAL '56 days'),
('50000000-0000-0000-0000-00000000000d', '40000000-0000-0000-0000-000000000005', 'break 语句', 'break 语句的作用是？', 'essay', 'medium', 15, 1, 600, 'null', '{"keywords": ["跳出", "循环", "终止", "退出"], "min_length": 20}', 3, NOW() - INTERVAL '56 days');

-- 第 6 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000000e', '40000000-0000-0000-0000-000000000006', '函数定义', '定义函数使用什么关键字？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "def", "alternatives": ["DEF"]}', 1, NOW() - INTERVAL '55 days'),
('50000000-0000-0000-0000-00000000000f', '40000000-0000-0000-0000-000000000006', '编写函数', '编写一个函数，接收两个参数并返回它们的和', 'coding', 'medium', 25, 5, 400, 'null', '{"expected_function": "add", "test_cases": [{"input": "2,3", "expected": "5"}, {"input": "-1,1", "expected": "0"}]}', 2, NOW() - INTERVAL '55 days'),
('50000000-0000-0000-0000-000000000010', '40000000-0000-0000-0000-000000000006', '默认参数', '函数可以有默认参数值。', 'true_false', 'easy', 5, 2, 30, 'null', '{"answer": true}', 3, NOW() - INTERVAL '55 days');

-- 第 7 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000011', '40000000-0000-0000-0000-000000000007', '列表索引', '列表 list = [1,2,3,4,5]，list[-1] 的值是？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "1", "is_correct": false}, {"text": "4", "is_correct": false}, {"text": "5", "is_correct": true}, {"text": "报错", "is_correct": false}]', '{"correct_option": 2}', 1, NOW() - INTERVAL '54 days'),
('50000000-0000-0000-0000-000000000012', '40000000-0000-0000-0000-000000000007', '列表推导式', '使用列表推导式生成 1-10 的平方', 'coding', 'medium', 25, 5, 400, 'null', '{"expected_output": "[1, 4, 9, 16, 25, 36, 49, 64, 81, 100]", "test_cases": [{"input": "", "expected": "[1, 4, 9, 16, 25, 36, 49, 64, 81, 100]"}]}', 2, NOW() - INTERVAL '54 days'),
('50000000-0000-0000-0000-000000000013', '40000000-0000-0000-0000-000000000007', '元组可变性', '元组创建后可以修改其中的元素。', 'true_false', 'easy', 5, 2, 30, 'null', '{"answer": false}', 3, NOW() - INTERVAL '54 days');

-- 第 8 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000014', '40000000-0000-0000-0000-000000000008', '字典创建', '以下哪个是正确的字典创建方式？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "dict = {key: value}", "is_correct": true}, {"text": "dict = [key: value]", "is_correct": false}, {"text": "dict = (key: value)", "is_correct": false}, {"text": "dict = <key: value>", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '53 days'),
('50000000-0000-0000-0000-000000000015', '40000000-0000-0000-0000-000000000008', '字典操作', '编写代码统计字符串中每个字符出现的次数', 'coding', 'hard', 30, 5, 600, 'null', '{"expected_function": "count_chars", "test_cases": [{"input": "hello", "expected": "{''h'': 1, ''e'': 1, ''l'': 2, ''o'': 1}"}]}', 2, NOW() - INTERVAL '53 days'),
('50000000-0000-0000-0000-000000000016', '40000000-0000-0000-0000-000000000008', '集合去重', '集合的主要特点是？', 'essay', 'medium', 15, 1, 600, 'null', '{"keywords": ["无序", "不重复", "唯一"], "min_length": 20}', 3, NOW() - INTERVAL '53 days');

-- 第 9 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000017', '40000000-0000-0000-0000-000000000009', '类定义', '定义类使用什么关键字？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "class", "alternatives": ["CLASS"]}', 1, NOW() - INTERVAL '52 days'),
('50000000-0000-0000-0000-000000000018', '40000000-0000-0000-0000-000000000009', '创建类', '创建一个 Person 类，包含 name 属性和 say_hello 方法', 'coding', 'hard', 35, 5, 600, 'null', '{"expected_class": "Person", "test_methods": ["say_hello"]}', 2, NOW() - INTERVAL '52 days'),
('50000000-0000-0000-0000-000000000019', '40000000-0000-0000-0000-000000000009', '继承概念', 'Python 支持多重继承。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": true}', 3, NOW() - INTERVAL '52 days');

-- 第 10 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000001a', '40000000-0000-0000-0000-00000000000a', '文件打开模式', '以只读方式打开文件使用什么模式？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "r", "is_correct": true}, {"text": "w", "is_correct": false}, {"text": "a", "is_correct": false}, {"text": "x", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '51 days'),
('50000000-0000-0000-0000-00000000001b', '40000000-0000-0000-0000-00000000000a', '文件读写', '编写代码读取文件内容并统计行数', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_function": "count_lines", "test_cases": [{"input": "test.txt", "expected": "integer"}]}', 2, NOW() - INTERVAL '51 days'),
('50000000-0000-0000-0000-00000000001c', '40000000-0000-0000-0000-00000000000a', '异常处理', 'try-except 结构中，无论是否发生异常都会执行的块是？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "finally", "alternatives": ["FINALLY"]}', 3, NOW() - INTERVAL '51 days');

-- ============ 课程 2: 机器学习基础 练习题 ============

-- 第 1 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000021', '40000000-0000-0000-0000-00000000000b', '机器学习定义', '机器学习是人工智能的一个子领域。', 'true_false', 'easy', 5, 2, 30, 'null', '{"answer": true}', 1, NOW() - INTERVAL '55 days'),
('50000000-0000-0000-0000-000000000022', '40000000-0000-0000-0000-00000000000b', '学习类型', '以下哪些属于机器学习的主要类型？', 'multiple_choice', 'easy', 15, 3, 90, '[{"text": "监督学习", "is_correct": true}, {"text": "无监督学习", "is_correct": true}, {"text": "强化学习", "is_correct": true}, {"text": "手动学习", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 2, NOW() - INTERVAL '55 days');

-- 第 2 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000023', '40000000-0000-0000-0000-00000000000c', 'NumPy 数组', '创建 NumPy 数组使用哪个函数？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "np.array()", "is_correct": true}, {"text": "np.create()", "is_correct": false}, {"text": "np.new()", "is_correct": false}, {"text": "np.list()", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '54 days'),
('50000000-0000-0000-0000-000000000024', '40000000-0000-0000-0000-00000000000c', 'Pandas DataFrame', '使用 Pandas 创建一个包含两列的 DataFrame', 'coding', 'medium', 25, 5, 400, 'null', '{"expected_class": "DataFrame", "columns": ["A", "B"]}', 2, NOW() - INTERVAL '54 days'),
('50000000-0000-0000-0000-000000000025', '40000000-0000-0000-0000-00000000000c', 'Matplotlib 绘图', 'Matplotlib 中用于绘制折线图的函数是？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "plot", "alternatives": ["plt.plot", "pyplot.plot"]}', 3, NOW() - INTERVAL '54 days');

-- 第 3 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000026', '40000000-0000-0000-0000-00000000000d', '缺失值处理', '以下哪种方法可以处理缺失值？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "删除含缺失值的行", "is_correct": true}, {"text": "用均值填充", "is_correct": true}, {"text": "用中位数填充", "is_correct": true}, {"text": "忽略不管", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 1, NOW() - INTERVAL '53 days'),
('50000000-0000-0000-0000-000000000027', '40000000-0000-0000-0000-00000000000d', '特征缩放', '标准化的公式是 (x - mean) / std。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": true}', 2, NOW() - INTERVAL '53 days'),
('50000000-0000-0000-0000-000000000028', '40000000-0000-0000-0000-00000000000d', '独热编码', '编写代码对分类变量进行独热编码', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_function": "one_hot_encode", "library": "pandas"}', 3, NOW() - INTERVAL '53 days');

-- 第 4 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000029', '40000000-0000-0000-0000-00000000000e', '线性回归假设', '线性回归假设自变量和因变量之间存在什么关系？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "线性", "alternatives": ["直线", "line"]}', 1, NOW() - INTERVAL '52 days'),
('50000000-0000-0000-0000-00000000002a', '40000000-0000-0000-0000-00000000000e', '实现线性回归', '使用 sklearn 实现线性回归模型', 'coding', 'hard', 35, 5, 600, 'null', '{"expected_class": "LinearRegression", "library": "sklearn"}', 2, NOW() - INTERVAL '52 days'),
('50000000-0000-0000-0000-00000000002b', '40000000-0000-0000-0000-00000000000e', 'R 平方值', 'R 平方值越接近 1，模型拟合效果越好。', 'true_false', 'easy', 10, 2, 60, 'null', '{"answer": true}', 3, NOW() - INTERVAL '52 days');

-- 第 5 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000002c', '40000000-0000-0000-0000-00000000000f', '逻辑回归输出', '逻辑回归的输出是什么类型的值？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "0 到 1 之间的概率", "is_correct": true}, {"text": "任意实数", "is_correct": false}, {"text": "整数", "is_correct": false}, {"text": "类别标签", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '51 days'),
('50000000-0000-0000-0000-00000000002d', '40000000-0000-0000-0000-00000000000f', 'Sigmoid 函数', '逻辑回归使用的激活函数是？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "sigmoid", "alternatives": ["Sigmoid", "逻辑函数", "logistic"]}', 2, NOW() - INTERVAL '51 days'),
('50000000-0000-0000-0000-00000000002e', '40000000-0000-0000-0000-00000000000f', '分类阈值', '描述如何选择合适的分类阈值', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["精确率", "召回率", "F1", "业务需求"], "min_length": 50}', 3, NOW() - INTERVAL '51 days');

-- 第 6 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000002f', '40000000-0000-0000-0000-000000000010', '决策树分裂', '决策树常用的分裂准则有？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "信息增益", "is_correct": true}, {"text": "基尼系数", "is_correct": true}, {"text": "均方误差", "is_correct": true}, {"text": "随机选择", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 1, NOW() - INTERVAL '50 days'),
('50000000-0000-0000-0000-000000000030', '40000000-0000-0000-0000-000000000010', '随机森林原理', '随机森林通过什么方式提高模型性能？', 'essay', 'hard', 25, 1, 900, 'null', '{"keywords": ["集成", "Bagging", "随机", "多棵树", "投票"], "min_length": 50}', 2, NOW() - INTERVAL '50 days'),
('50000000-0000-0000-0000-000000000031', '40000000-0000-0000-0000-000000000010', '实现随机森林', '使用 sklearn 创建随机森林分类器', 'coding', 'hard', 30, 5, 600, 'null', '{"expected_class": "RandomForestClassifier", "library": "sklearn"}', 3, NOW() - INTERVAL '50 days');

-- 第 7 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000032', '40000000-0000-0000-0000-000000000011', 'SVM 核心思想', 'SVM 的核心思想是找到最大间隔的超平面。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": true}', 1, NOW() - INTERVAL '49 days'),
('50000000-0000-0000-0000-000000000033', '40000000-0000-0000-0000-000000000011', '核函数作用', '核函数的作用是？', 'multiple_choice', 'hard', 20, 3, 120, '[{"text": "将数据映射到高维空间", "is_correct": true}, {"text": "降低数据维度", "is_correct": false}, {"text": "删除异常值", "is_correct": false}, {"text": "增加训练数据", "is_correct": false}]', '{"correct_option": 0}', 2, NOW() - INTERVAL '49 days'),
('50000000-0000-0000-0000-000000000034', '40000000-0000-0000-0000-000000000011', '常用核函数', '列举三种常用的核函数', 'fill_blank', 'hard', 20, 3, 180, 'null', '{"answer": "linear,rbf,poly", "alternatives": ["线性，径向基，多项式"]}', 3, NOW() - INTERVAL '49 days');

-- 第 8 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000035', '40000000-0000-0000-0000-000000000012', 'K-Means 输入', 'K-Means 算法需要预先指定什么参数？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "聚类数量 K", "is_correct": true}, {"text": "学习率", "is_correct": false}, {"text": "迭代次数", "is_correct": false}, {"text": "正则化参数", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '48 days'),
('50000000-0000-0000-0000-000000000036', '40000000-0000-0000-0000-000000000012', '实现 K-Means', '使用 sklearn 实现 K-Means 聚类', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_class": "KMeans", "library": "sklearn"}', 2, NOW() - INTERVAL '48 days'),
('50000000-0000-0000-0000-000000000037', '40000000-0000-0000-0000-000000000012', '肘部法则', '肘部法则用于确定最优的 K 值。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": true}', 3, NOW() - INTERVAL '48 days');

-- 第 9 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000038', '40000000-0000-0000-0000-000000000013', '准确率局限', '在类别不平衡的情况下，准确率是可靠的评估指标。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": false}', 1, NOW() - INTERVAL '47 days'),
('50000000-0000-0000-0000-000000000039', '40000000-0000-0000-0000-000000000013', '交叉验证', 'K 折交叉验证中，数据被分成几份？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "K", "alternatives": ["k", "若干"]}', 2, NOW() - INTERVAL '47 days'),
('50000000-0000-0000-0000-00000000003a', '40000000-0000-0000-0000-000000000013', '网格搜索', '使用网格搜索优化模型超参数', 'coding', 'hard', 30, 5, 600, 'null', '{"expected_class": "GridSearchCV", "library": "sklearn"}', 3, NOW() - INTERVAL '47 days');

-- 第 10 章练习
INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-00000000003b', '40000000-0000-0000-0000-000000000014', '特征工程', '在房价预测中，哪些特征可能重要？', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["面积", "位置", "房龄", "楼层", "朝向", "交通"], "min_length": 50}', 1, NOW() - INTERVAL '46 days'),
('50000000-0000-0000-0000-00000000003c', '40000000-0000-0000-0000-000000000014', '模型选择', '对于房价预测问题，最适合的算法是？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "线性回归", "is_correct": true}, {"text": "K-Means", "is_correct": false}, {"text": "逻辑回归", "is_correct": false}, {"text": "SVM 分类", "is_correct": false}]', '{"correct_option": 0}', 2, NOW() - INTERVAL '46 days'),
('50000000-0000-0000-0000-00000000003d', '40000000-0000-0000-0000-000000000014', '完整项目', '实现完整的房价预测流程', 'coding', 'hard', 50, 5, 1200, 'null', '{"expected_steps": ["load_data", "preprocess", "train", "evaluate", "predict"]}', 3, NOW() - INTERVAL '46 days');

-- ============ 课程 3: Web 开发实战 练习题 ============
-- (简化版本，实际应包含 10+ 题，这里为节省空间省略部分)

INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000041', '40000000-0000-0000-0000-000000000015', 'HTTP 方法', 'GET 和 POST 的主要区别是什么？', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["参数位置", "安全性", "幂等性", "缓存"], "min_length": 50}', 1, NOW() - INTERVAL '50 days'),
('50000000-0000-0000-0000-000000000042', '40000000-0000-0000-0000-000000000015', '状态码', 'HTTP 状态码 200 表示？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "请求成功", "is_correct": true}, {"text": "未找到", "is_correct": false}, {"text": "服务器错误", "is_correct": false}, {"text": "重定向", "is_correct": false}]', '{"correct_option": 0}', 2, NOW() - INTERVAL '50 days'),
('50000000-0000-0000-0000-000000000043', '40000000-0000-0000-0000-000000000016', 'HTML 语义化', '以下哪个是语义化标签？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "<article>", "is_correct": true}, {"text": "<div>", "is_correct": false}, {"text": "<span>", "is_correct": false}, {"text": "<i>", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '49 days'),
('50000000-0000-0000-0000-000000000044', '40000000-0000-0000-0000-000000000016', '创建表单', '创建一个包含用户名和密码输入框的 HTML 表单', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_tags": ["form", "input", "button"]}', 2, NOW() - INTERVAL '49 days'),
('50000000-0000-0000-0000-000000000045', '40000000-0000-0000-0000-000000000017', 'CSS 选择器', '选择 id 为"main"的元素使用什么选择器？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "#main", "alternatives": ["#MAIN"]}', 1, NOW() - INTERVAL '48 days'),
('50000000-0000-0000-0000-000000000046', '40000000-0000-0000-0000-000000000017', 'Flexbox 布局', 'Flexbox 中 justify-content 的作用是？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "主轴对齐", "is_correct": true}, {"text": "交叉轴对齐", "is_correct": false}, {"text": "设置 flex 增长", "is_correct": false}, {"text": "设置 flex 收缩", "is_correct": false}]', '{"correct_option": 0}', 2, NOW() - INTERVAL '48 days'),
('50000000-0000-0000-0000-000000000047', '40000000-0000-0000-0000-000000000018', 'DOM 操作', '使用 JavaScript 获取 id 为"demo"的元素', 'coding', 'easy', 15, 5, 300, 'null', '{"expected_method": "getElementById", "expected_selector": "demo"}', 1, NOW() - INTERVAL '47 days'),
('50000000-0000-0000-0000-000000000048', '40000000-0000-0000-0000-000000000018', '事件监听', 'addEventListener 可以绑定多个事件处理函数。', 'true_false', 'easy', 10, 2, 60, 'null', '{"answer": true}', 2, NOW() - INTERVAL '47 days'),
('50000000-0000-0000-0000-000000000049', '40000000-0000-0000-0000-000000000019', 'Promise 状态', 'Promise 有哪几种状态？', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["pending", "fulfilled", "resolved", "rejected"], "min_length": 30}', 1, NOW() - INTERVAL '46 days'),
('50000000-0000-0000-0000-00000000004a', '40000000-0000-0000-0000-000000000019', 'Fetch API', '使用 Fetch API 发送 GET 请求', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_method": "fetch", "expected_params": ["url"]}', 2, NOW() - INTERVAL '46 days'),
('50000000-0000-0000-0000-00000000004b', '40000000-0000-0000-0000-00000000001a', 'Express 路由', '在 Express 中定义 GET 路由使用？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "app.get", "alternatives": ["router.get", "app.get()"]}', 1, NOW() - INTERVAL '45 days'),
('50000000-0000-0000-0000-00000000004c', '40000000-0000-0000-0000-00000000001a', '中间件概念', 'Express 中间件可以做什么？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "修改请求对象", "is_correct": true}, {"text": "修改响应对象", "is_correct": true}, {"text": "终止请求响应循环", "is_correct": true}, {"text": "调用下一个中间件", "is_correct": true}]', '{"correct_options": [0, 1, 2, 3]}', 2, NOW() - INTERVAL '45 days'),
('50000000-0000-0000-0000-00000000004d', '40000000-0000-0000-0000-00000000001b', 'MongoDB 操作', 'MongoDB 中插入文档使用什么方法？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "insertOne/insertMany", "is_correct": true}, {"text": "addOne", "is_correct": false}, {"text": "create", "is_correct": false}, {"text": "save", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '44 days'),
('50000000-0000-0000-0000-00000000004e', '40000000-0000-0000-0000-00000000001b', '查询文档', '使用 MongoDB 查询 age 大于 18 的文档', 'coding', 'medium', 25, 5, 500, 'null', '{"expected_method": "find", "expected_query": {"age": {"$gt": 18}}}', 2, NOW() - INTERVAL '44 days'),
('50000000-0000-0000-0000-00000000004f', '40000000-0000-0000-0000-00000000001c', 'JWT 结构', 'JWT 由哪三部分组成？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "header.payload.signature", "alternatives": ["头部。载荷。签名", "Header.Payload.Signature"]}', 1, NOW() - INTERVAL '43 days'),
('50000000-0000-0000-0000-000000000050', '40000000-0000-0000-0000-00000000001c', '实现登录', '实现用户登录接口，返回 JWT Token', 'coding', 'hard', 40, 5, 800, 'null', '{"expected_features": ["password_hash", "token_generation", "expiration"]}', 2, NOW() - INTERVAL '43 days');

-- ============ 课程 4: 数据结构与算法 练习题 ============

INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000051', '40000000-0000-0000-0000-00000000001f', '时间复杂度', 'O(1) 表示什么时间复杂度？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "常数时间", "is_correct": true}, {"text": "线性时间", "is_correct": false}, {"text": "对数时间", "is_correct": false}, {"text": "平方时间", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '45 days'),
('50000000-0000-0000-0000-000000000052', '40000000-0000-0000-0000-00000000001f', '复杂度比较', '按时间复杂度从小到大排序：O(n), O(1), O(n²), O(log n)', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["O(1)", "O(log n)", "O(n)", "O(n²)"], "min_length": 20}', 2, NOW() - INTERVAL '45 days'),
('50000000-0000-0000-0000-000000000053', '40000000-0000-0000-0000-000000000020', '数组访问', '数组的随机访问时间复杂度是？', 'fill_blank', 'easy', 10, 3, 60, 'null', '{"answer": "O(1)", "alternatives": ["常数", "constant"]}', 1, NOW() - INTERVAL '44 days'),
('50000000-0000-0000-0000-000000000054', '40000000-0000-0000-0000-000000000020', '两数之和', '实现两数之和算法，返回数组中两个数的索引', 'coding', 'medium', 30, 5, 600, 'null', '{"expected_function": "twoSum", "test_cases": [{"input": "[2,7,11,15],9", "expected": "[0,1]"}]}', 2, NOW() - INTERVAL '44 days'),
('50000000-0000-0000-0000-000000000055', '40000000-0000-0000-0000-000000000021', '链表反转', '反转单链表的时间复杂度是？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "O(n)", "is_correct": true}, {"text": "O(1)", "is_correct": false}, {"text": "O(n²)", "is_correct": false}, {"text": "O(log n)", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '43 days'),
('50000000-0000-0000-0000-000000000056', '40000000-0000-0000-0000-000000000021', '实现链表', '实现一个单链表类，包含 insert 和 delete 方法', 'coding', 'hard', 40, 5, 900, 'null', '{"expected_class": "LinkedList", "expected_methods": ["insert", "delete", "search"]}', 2, NOW() - INTERVAL '43 days'),
('50000000-0000-0000-0000-000000000057', '40000000-0000-0000-0000-000000000022', '栈特性', '栈的特点是后进先出 (LIFO)。', 'true_false', 'easy', 10, 2, 60, 'null', '{"answer": true}', 1, NOW() - INTERVAL '42 days'),
('50000000-0000-0000-0000-000000000058', '40000000-0000-0000-0000-000000000022', '括号匹配', '使用栈实现括号匹配检查', 'coding', 'medium', 30, 5, 600, 'null', '{"expected_function": "isValid", "test_cases": [{"input": "()[]{}", "expected": "true"}, {"input": "([)]", "expected": "false"}]}', 2, NOW() - INTERVAL '42 days'),
('50000000-0000-0000-0000-000000000059', '40000000-0000-0000-0000-000000000023', '哈希冲突', '解决哈希冲突的方法有？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "链地址法", "is_correct": true}, {"text": "开放寻址法", "is_correct": true}, {"text": "再哈希法", "is_correct": true}, {"text": "删除冲突元素", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 1, NOW() - INTERVAL '41 days'),
('50000000-0000-0000-0000-00000000005a', '40000000-0000-0000-0000-000000000023', '实现哈希表', '实现一个简单的哈希表，支持 put 和 get 操作', 'coding', 'hard', 40, 5, 900, 'null', '{"expected_class": "HashMap", "expected_methods": ["put", "get", "remove"]}', 2, NOW() - INTERVAL '41 days'),
('50000000-0000-0000-0000-00000000005b', '40000000-0000-0000-0000-000000000024', '二叉树遍历', '二叉树的前序遍历顺序是？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "根左右", "is_correct": true}, {"text": "左根右", "is_correct": false}, {"text": "左右根", "is_correct": false}, {"text": "右左根", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '40 days'),
('50000000-0000-0000-0000-00000000005c', '40000000-0000-0000-0000-000000000024', '实现 BST', '实现二叉搜索树的插入和查找操作', 'coding', 'hard', 40, 5, 900, 'null', '{"expected_class": "BST", "expected_methods": ["insert", "search", "delete"]}', 2, NOW() - INTERVAL '40 days');

-- ============ 课程 5: 深度学习入门 练习题 ============

INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, time_limit, options, expected_answer, order_index, created_at) VALUES
('50000000-0000-0000-0000-000000000061', '40000000-0000-0000-0000-000000000029', '深度学习应用', '深度学习已成功应用于哪些领域？', 'multiple_choice', 'easy', 15, 3, 90, '[{"text": "图像识别", "is_correct": true}, {"text": "语音识别", "is_correct": true}, {"text": "自然语言处理", "is_correct": true}, {"text": "游戏 AI", "is_correct": true}]', '{"correct_options": [0, 1, 2, 3]}', 1, NOW() - INTERVAL '40 days'),
('50000000-0000-0000-0000-000000000062', '40000000-0000-0000-0000-000000000029', '神经网络优势', '相比传统机器学习，深度学习的优势是什么？', 'essay', 'medium', 20, 1, 600, 'null', '{"keywords": ["自动特征", "端到端", "大数据", "表达能力"], "min_length": 50}', 2, NOW() - INTERVAL '40 days'),
('50000000-0000-0000-0000-000000000063', '40000000-0000-0000-0000-00000000002a', '激活函数', '以下哪些是常用的激活函数？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "ReLU", "is_correct": true}, {"text": "Sigmoid", "is_correct": true}, {"text": "Tanh", "is_correct": true}, {"text": "Linear", "is_correct": false}]', '{"correct_options": [0, 1, 2]}', 1, NOW() - INTERVAL '39 days'),
('50000000-0000-0000-0000-000000000064', '40000000-0000-0000-0000-00000000002a', '反向传播', '反向传播算法的作用是？', 'fill_blank', 'medium', 15, 3, 90, 'null', '{"answer": "计算梯度", "alternatives": ["梯度计算", "更新权重"]}', 2, NOW() - INTERVAL '39 days'),
('50000000-0000-0000-0000-000000000065', '40000000-0000-0000-0000-00000000002b', 'TensorFlow 张量', 'TensorFlow 中创建常量张量使用？', 'multiple_choice', 'easy', 10, 3, 60, '[{"text": "tf.constant", "is_correct": true}, {"text": "tf.create", "is_correct": false}, {"text": "tf.new", "is_correct": false}, {"text": "tf.tensor", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '38 days'),
('50000000-0000-0000-0000-000000000066', '40000000-0000-0000-0000-00000000002b', '构建模型', '使用 TensorFlow Keras 构建一个简单的全连接网络', 'coding', 'hard', 40, 5, 900, 'null', '{"expected_api": "tf.keras", "expected_layers": ["Dense", "Activation"]}', 2, NOW() - INTERVAL '38 days'),
('50000000-0000-0000-0000-000000000067', '40000000-0000-0000-0000-00000000002c', 'PyTorch 特点', 'PyTorch 使用静态计算图。', 'true_false', 'medium', 10, 2, 60, 'null', '{"answer": false}', 1, NOW() - INTERVAL '37 days'),
('50000000-0000-0000-0000-000000000068', '40000000-0000-0000-0000-00000000002c', '定义网络', '使用 PyTorch 定义一个神经网络类', 'coding', 'hard', 40, 5, 900, 'null', '{"expected_base": "nn.Module", "expected_methods": ["__init__", "forward"]}', 2, NOW() - INTERVAL '37 days'),
('50000000-0000-0000-0000-000000000069', '40000000-0000-0000-0000-00000000002d', '卷积操作', '卷积层的主要作用是？', 'multiple_choice', 'medium', 15, 3, 90, '[{"text": "提取局部特征", "is_correct": true}, {"text": "降低维度", "is_correct": false}, {"text": "增加非线性", "is_correct": false}, {"text": "正则化", "is_correct": false}]', '{"correct_option": 0}', 1, NOW() - INTERVAL '36 days'),
('50000000-0000-0000-0000-00000000006a', '40000000-0000-0000-0000-00000000002d', '实现 CNN', '使用 PyTorch 实现一个简单的 CNN 用于图像分类', 'coding', 'hard', 50, 5, 1200, 'null', '{"expected_layers": ["Conv2d", "ReLU", "MaxPool2d", "Linear"]}', 2, NOW() - INTERVAL '36 days'),
('50000000-0000-0000-0000-00000000006b', '40000000-0000-0000-0000-00000000002e', 'RNN 应用', 'RNN 适合处理序列数据。', 'true_false', 'easy', 10, 2, 60, 'null', '{"answer": true}', 1, NOW() - INTERVAL '35 days'),
('50000000-0000-0000-0000-00000000006c', '40000000-0000-0000-0000-00000000002e', 'LSTM 优势', 'LSTM 相比普通 RNN 的优势是什么？', 'essay', 'hard', 25, 1, 900, 'null', '{"keywords": ["长依赖", "梯度消失", "门控", "记忆单元"], "min_length": 50}', 2, NOW() - INTERVAL '35 days');
