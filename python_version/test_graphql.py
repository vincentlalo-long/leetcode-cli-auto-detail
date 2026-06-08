import requests, json

query = """query getQuestionDetail($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    similarQuestions
  }
}"""
variables = {'titleSlug': 'two-sum'}
response = requests.post('https://leetcode.com/graphql', json={'query': query, 'variables': variables})
print(response.json())