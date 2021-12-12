package collectors

var (
	aggregationsFile string = "\n" +
		"===========\n" +
		"Aggregation\n" +
		"===========\n" +
		"\n" +
		".. default-domain:: mongodb\n" +
		"\n" +
		".. contents:: On this page\n" +
		":local:\n" +
		":backlinks: none\n" +
		":depth: 2\n" +
		":class: singlecol\n" +
		"\n" +
		".. _nodejs-aggregation-overview:\n" +
		"\n" +
		"Overview\n" +
		"--------\n" +
		"\n" +
		"In this guide, you can learn how to use **aggregation operations** in\n" +
		"the MongoDB Node.js driver.\n" +
		"\n" +
		"Aggregation operations are expressions you can use to produce reduced\n" +
		"and summarized results in MongoDB. MongoDB's aggregation framework\n" +
		"allows you to create a pipeline that consists of one or more stages,\n" +
		"each of which performs a specific operation on your data.\n" +
		"\n" +
		"You can think of the aggregation framework as similar to an automobile factory.\n" +
		"Automobile manufacturing requires the use of assembly stations organized\n" +
		"into assembly lines. Each station has specialized tools, such as\n" +
		"drills and welders. The factory transforms and\n" +
		"assembles the initial parts and materials into finished products.\n" +
		"\n" +
		"The **aggregation pipeline** is the assembly line, **aggregation stages** are the assembly stations, and\n" +
		"**operator expressions** are the specialized tools.\n" +
		"\n" +
		"Aggregation vs. Query Operations\n" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" +
		"\n" +
		"Using query operations, such as the ``find()`` method, you can perform the following actions:\n" +
		"\n" +
		"- Select *which documents* to return.\n" +
		"- Select *which fields* to return.\n" +
		"- Sort the results.\n" +
		"\n" +
		"Using aggregation operations, you can perform the following actions:\n" +
		"\n" +
		"- Perform all query operations.\n" +
		"- Rename fields.\n" +
		"- Calculate fields.\n" +
		"- Summarize data.\n" +
		"- Group values.\n" +
		"\n" +
		"Aggregation operations have some :manual:`limitations </core/aggregation-pipeline-limits/>`:\n" +
		"\n" +
		"- Returned documents must not violate the :manual:`BSON-document size limit </reference/limits/#mongodb-limit-BSON-Document-Size>`\n" +
		"of 16 megabytes.\n" +
		"\n" +
		"- Pipeline stages have a memory limit of 100 megabytes by default. If necessary, you may exceed this limit by setting the ``allowDiskUse``\n" +
		"property of ``AggregateOptions`` to ``true``. See the\n" +
		"`AggregateOptions API documentation <{+api+}/interfaces/AggregateOptions.html>`__\n" +
		"for more details.\n" +
		"\n" +
		".. important:: ``$graphLookup`` exception\n" +
		"\n" +
		"The :manual:`$graphLookup\n" +
		"</reference/operator/aggregation/graphLookup/>` stage has a strict\n" +
		"memory limit of 100 megabytes and will ignore ``allowDiskUse``.\n" +
		"\n" +
		"Useful References\n" +
		"~~~~~~~~~~~~~~~~~\n" +
		"\n" +
		"- :manual:`Expression operators </reference/operator/aggregation/>`\n" +
		"- :manual:`Aggregation pipeline </core/aggregation-pipeline/>`\n" +
		"- :manual:`Aggregation stages </meta/aggregation-quick-reference/#stages>`\n" +
		"- :manual:`Operator expressions </meta/aggregation-quick-reference/#operator-expressions>`\n" +
		"\n" +
		"Runnable Examples\n" +
		"-----------------\n" +
		"\n" +
		"The example uses sample data about restaurants. The following code\n" +
		"inserts data into the ``restaurants`` collection of the ``aggregation``\n" +
		"database:\n" +
		"\n" +
		".. literalinclude:: /code-snippets/aggregation/agg.js\n" +
		":start-after: begin data insertion\n" +
		":end-before: end data insertion\n" +
		":language: javascript\n" +
		":dedent:\n" +
		"\n" +
		".. tip::\n" +
		"\n" +
		"For more information on connecting to your MongoDB deployment, see the :doc:`Connection Guide </fundamentals/connection>`.\n" +
		"\n" +
		"Aggregation Example\n" +
		"~~~~~~~~~~~~~~~~~~~\n" +
		"\n" +
		"To perform an aggregation, pass a list of aggregation stages to the\n" +
		"``collection.aggregate()`` method.\n" +
		"\n" +
		"In the example, the aggregation pipeline uses the following aggregation stages:\n" +
		"\n" +
		"- A :manual:`$match </reference/operator/aggregation/match/>` stage to filter for documents whose\n" +
		"``categories`` array field contains the element ``Bakery``.\n" +
		"\n" +
		"- A :manual:`$group </reference/operator/aggregation/group/>` stage to group the matching documents by the ``stars``\n" +
		"field, accumulating a count of documents for each distinct value of ``stars``.\n" +
		"\n" +
		".. literalinclude:: /code-snippets/aggregation/agg.js\n" +
		":start-after: begin aggregation\n" +
		":end-before: end aggregation\n" +
		":language: javascript\n" +
		":dedent:\n" +
		"\n" +
		"This example should produce the following output:\n" +
		"\n" +
		".. code-block:: json\n" +
		":copyable: false\n" +
		"\n" +
		"{ _id: 4, count: 2 }\n" +
		"{ _id: 3, count: 1 }\n" +
		"{ _id: 5, count: 1 }\n" +
		"\n" +
		"For more information, see the `aggregate() API documentation <{+api+}/classes/Collection.html#aggregate>`__.\n" +
		"\n" +
		"Additional Aggregation Examples\n" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" +
		"You can find another aggregation framework example `in this MongoDB Blog\n" +
		"post <https://www.mongodb.com/blog/post/quick-start-nodejs--mongodb--how-to-analyze-data-using-the-aggregation-framework>`_.\n"

	indexFile string = "" +
		"===================\n" +
		"MongoDB Node Driver\n" +
		"===================\n" +
		"\n" +
		".. default-domain:: mongodb\n" +
		"\n" +
		".. include:: /includes/unicode-checkmark.rst\n" +
		"\n" +
		".. toctree::\n" +
		":titlesonly:\n" +
		":maxdepth: 1\n" +
		"\n" +
		"/quick-start\n" +
		"/usage-examples\n" +
		"/fundamentals\n" +
		"API Documentation <{+api+}>\n" +
		"/faq\n" +
		"/issues-and-help\n" +
		"/compatibility\n" +
		"/whats-new\n" +
		"Release Notes <https://github.com/mongodb/node-mongodb-native/releases/>\n" +
		"View the Source <https://github.com/mongodb/node-mongodb-native/>\n" +
		"\n" +
		"Introduction\n" +
		"------------\n" +
		"\n" +
		"Welcome to the MongoDB Node.js driver documentation. \n" +
		"\n" +
		"Connect your Node.js applications to MongoDB and work with your data\n" +
		"using the official Node.js driver. The driver features an asynchronous\n" +
		"API that you can use to access method return values through Promises or\n" +
		"specify callbacks to access them when communicating with MongoDB.\n" +
		"\n" +
		"On this page, you can find descriptions of each section of the driver\n" +
		"documentation and how to learn more about the Node.js driver.\n" +
		"\n" +
		"Quick Start\n" +
		"-----------\n" +
		"\n" +
		"Learn how to establish a connection to MongoDB Atlas and begin\n" +
		"working with data in the :doc:`Quick Start </quick-start>` section.\n" +
		"\n" +
		"Usage Examples\n" +
		"--------------\n" +
		"\n" +
		"For fully runnable code snippets and explanations for common\n" +
		"methods, see the :doc:`Usage Examples </usage-examples>` section.\n" +
		"\n" +
		"Fundamentals\n" +
		"------------\n" +
		"\n" +
		".. include:: /includes/fundamentals-sections.rst\n" +
		"\n" +
		"API\n" +
		"---\n" +
		"\n" +
		"For detailed information about classes and methods in the MongoDB\n" +
		"Node.js driver, see the `MongoDB Node.js driver API documentation\n" +
		"<{+api+}>`__ .\n" +
		"\n" +
		"FAQ\n" +
		"---\n" +
		"\n" +
		"For answers to commonly asked questions about the MongoDB\n" +
		"Node.js Driver, see the :doc:`Frequently Asked Questions (FAQ) </faq>`\n" +
		"section.\n" +
		"\n" +
		"Issues & Help\n" +
		"-------------\n" +
		"\n" +
		"Learn how to report bugs, contribute to the driver, and find\n" +
		"additional resources for asking questions and receiving help in the\n" +
		":doc:`Issues & Help </issues-and-help>` section.\n" +
		"\n" +
		"Compatibility\n" +
		"-------------\n" +
		"\n" +
		"For the compatibility charts that show the recommended Node.js\n" +
		"Driver version for each MongoDB Server version, see the\n" +
		":doc:`Compatibility </compatibility>` section.\n" +
		"\n" +
		"What's New\n" +
		"----------\n" +
		"\n" +
		"For a list of new features and changes in each version, see the\n" +
		":doc:`What's New </whats-new>` section.\n" +
		"\n" +
		"Learn\n" +
		"-----\n" +
		"\n" +
		"Visit the Developer Hub and MongoDB University to learn more about the\n" +
		"MongoDB Node.js driver.\n" +
		"\n" +
		"Developer Hub\n" +
		"~~~~~~~~~~~~~\n" +
		"\n" +
		"The Developer Hub provides tutorials and social engagement for\n" +
		"developers.\n" +
		"\n" +
		"To learn how to use MongoDB features with the Node.js driver, see the\n" +
		"`How To's and Articles page <https://developer.mongodb.com/learn/?content=Articles&text=Node.js>`_.\n" +
		"\n" +
		"To ask questions and engage in discussions with fellow developers using\n" +
		"the Node.js driver, see the `forums page <https://developer.mongodb.com/community/forums/tag/node-js>`_.\n" +
		"\n" +
		"MongoDB University\n" +
		"~~~~~~~~~~~~~~~~~~\n" +
		"\n" +
		"MongoDB University provides free courses to teach everyone how to use MongoDB.\n" +
		"\n" +
		"Take the free online course taught by MongoDB instructors\n" +
		"`````````````````````````````````````````````````````````\n" +
		"\n" +
		".. list-table::\n" +
		"\n" +
		"* - .. cssclass:: bordered-figure\n" +
		".. figure:: /includes/figures/M220JS_hero.jpg\n" +
		":alt: M220JS course banner\n" +
		"\n" +
		"- `M220JS: MongoDB for JavaScript Developers <https://university.mongodb.com/courses/M220JS/about>`_\n" +
		"Learn the essentials of Node.js application development with\n" +
		"MongoDB.\n"
)
