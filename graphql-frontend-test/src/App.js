import React, { Component } from 'react';
import './App.css';
import ReactTable from 'react-table'
import 'react-table/react-table.css'
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
import { gql, graphql } from 'react-apollo';


const data = [{
  name: 'Tanner Linsley',
  age: 26,
  friend: {
    name: 'Jason Maurer',
    age: 23,
  }
}]

const columns = [{
  header: 'Name',
  accessor: 'name' // String-based value accessors!
}, {
  header: 'Age',
  accessor: 'age',
  render: props => <span className='number'>{props.value}</span> // Custom cell components!
}, {
  header: 'Friend Name',
  accessor: d => d.friend.name, // Custom value accessors!
  id: 'fname'
}, {
  header: props => <span>Friend Age</span>, // Custom header components!
  accessor: 'friend.age'
}]

function ProductList({ loading, products }) {
  if (loading) {
    return <div>Loading</div>;
  } else {
    return (
      <div className="App">
        <ReactTable
          data={data}
          columns={columns}
        />
        <ul>
          {products.map(product =>
              <li key={product.sku}>
                {product.title} by {' '}
                {product.author.firstName} {product.author.lastName} {' '}
                ({product.votes} votes)
              </li>
          )}
        </ul>
      </div>
    );
  }
}
const allProducts = gql`
query products {
    sku
    location
    instanceType
    operatingSystem
  }
`

const ProductListWithData = graphql(allProducts, {
  props: ({data: { loading, products }}) => ({
    loading,
      products,
  }),
})(ProductList);



class App extends Component {
  constructor(...args) {
    super(...args);
    const networkInterface = createNetworkInterface({
      uri: 'http://localhost:8080/graphql'
    })

    this.client = new ApolloClient({
      networkInterface: networkInterface
    });
  }

  render() {
    return (
      <ApolloProvider client={this.client}>
        <ProductListWithData />
      </ApolloProvider>
    );
  }
}



export default App;

