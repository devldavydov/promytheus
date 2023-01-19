package data

var AllMetricsEmptyResponse string = `
	<html>
		<body>
			<table border="1">
				<tr>
					<th>Metric Type</th>
					<th>Metric Name</th>
					<th>Metric Value</th>
				</tr>
				
			</table>
		</body>
	</html>
	`

var AllMetricsResponseWithData string = `
	<html>
		<body>
			<table border="1">
				<tr>
					<th>Metric Type</th>
					<th>Metric Name</th>
					<th>Metric Value</th>
				</tr>
				
				<tr>
					<td>counter</td>
					<td>aaa</td>
					<td>2</td>
				</tr>
				
				<tr>
					<td>counter</td>
					<td>zzz</td>
					<td>3</td>
				</tr>
				
				<tr>
					<td>gauge</td>
					<td>bar</td>
					<td>1.235</td>
				</tr>
				
				<tr>
					<td>gauge</td>
					<td>foo</td>
					<td>1.235</td>
				</tr>
				
			</table>
		</body>
	</html>
	`
