<html>
	<body>
		<h1>login</h1>
		<form action="/user/login" method="post">
			用户名:<input type="text" name="username"><br>
			密码:<input type="password" name="password"><br>
			<input type="submit" value="登录"><br>
			<a href="/user/signup">注册</a>
		</form>
		<div style="position:absolute;bottom:0">{{.content}}<div>
	</body>
</html>